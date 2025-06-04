package ui

import (
	"fmt"
	"os"
	"os/exec"
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

	// check if the command is a function call
	// if strings.HasPrefix(command, "kit ") {
	// 	command = fmt.Sprintf("%s %s", os.Getenv("MENACE_VENV_PATH"), command)
	// }
	

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	output, err := cmd.CombinedOutput()
	return string(output), err
}
