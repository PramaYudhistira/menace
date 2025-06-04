package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type DiffType int

// 3 types of file changes, Add, Delete, Modify
const (
	Add DiffType = iota
	Delete
	Modify
)

// Represents actions to do for each line
type LineDiff struct {
	Type       DiffType //0 for add
	LineIndex  int      //Which line we are trying to modify, 1-indexed
	OldContent string
	NewContent string
}

// Helper function to write diffs to file
func applyDiffsToFile(path string, diffs []LineDiff) error {
	raw, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	lines := []string{}
	if err == nil {
		lines = strings.Split(string(raw), "\n")
	}

	// --- Deletes (reverse order) ---
	for i := len(diffs) - 1; i >= 0; i-- {
		d := diffs[i]
		if d.Type == Delete {
			idx := d.LineIndex - 1 // convert to 0‑based
			if idx >= 0 && idx < len(lines) {
				lines = append(lines[:idx], lines[idx+1:]...)
			}
		}
	}

	// --- Adds & Modifies (ascending) ---
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].LineIndex < diffs[j].LineIndex
	})
	for _, d := range diffs {
		idx := d.LineIndex - 1 // convert to 0‑based
		switch d.Type {
		case Add:
			if idx < 0 {
				idx = 0
			}
			if idx > len(lines) {
				idx = len(lines)
			}
			lines = append(lines[:idx],
				append([]string{d.NewContent}, lines[idx:]...)...)
		case Modify:
			if idx >= 0 && idx < len(lines) {
				lines[idx] = d.NewContent
			}
		}
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

// Writes to an existing file
//
// If file doesn't exist, creates the file and writes to it
// Returns an error if anything happens, should be emitted to the model
func CreateAndApplyDiffs(path string, diffs []LineDiff) error {
	// make sure the file exists
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot create or open %s: %w", path, err)
	}
	f.Close()

	//  now apply any diffs
	return applyDiffsToFile(path, diffs)
}

// Returns string, error
//
// string is a zero indexed array of line numbers and text
func ReadFileWithLineNumbers(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	var sb strings.Builder
	for i, line := range lines {
		// change i → i+1
		sb.WriteString(fmt.Sprintf("%d: %s\n", i+1, line))
	}
	return sb.String(), nil
}

func FormatDiff(d LineDiff) string {
	switch d.Type {
	case Add:
		return fmt.Sprintf("  + [%d] %s", d.LineIndex, d.NewContent)
	case Delete:
		return fmt.Sprintf("  - [%d] %s", d.LineIndex, d.OldContent)
	case Modify:
		return fmt.Sprintf("  ~ [%d] %s → %s",
			d.LineIndex, d.OldContent, d.NewContent)
	default:
		return "  ? unknown diff"
	}
}

func PrintDiffs(diffs []LineDiff) {
	fmt.Println("Applying diffs:")
	for _, d := range diffs {
		fmt.Println(FormatDiff(d))
	}
}

// Runs a shell command
func runShellCommand(command string) (string, error) {
	// Check if the command is a kit command
	if strings.HasPrefix(command, "kit ") {
		venvPath := os.Getenv("MENACE_VENV_PATH")
		if venvPath == "" {
			return "", fmt.Errorf("MENACE_VENV_PATH not set")
		}

		// Get the kit executable path
		kitPath := filepath.Join(venvPath, "bin", "kit")
		if runtime.GOOS == "windows" {
			kitPath = filepath.Join(venvPath, "Scripts", "kit.exe")
		}

		// Split the command to get kit subcommand and arguments
		parts := strings.Fields(command)
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid kit command format")
		}

		// Create the command with kit path and remaining arguments
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command(kitPath, parts[1:]...)
		} else {
			cmd = exec.Command(kitPath, parts[1:]...)
		}

		// Set up environment
		cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s:%s",
			filepath.Join(venvPath, "bin"),
			os.Getenv("PATH")))

		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("kit command failed: %v\nOutput: %s", err, string(output))
		}
		return string(output), nil
	}

	// Handle regular shell commands
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	output, err := cmd.CombinedOutput()
	return string(output), err
}
