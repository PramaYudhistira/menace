package llmServer

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ModelFactory struct {
	Model string
}

// TODO: Refactor?
func (mf ModelFactory) DetectShell() string {
	goos := runtime.GOOS
	var shell string

	if goos == "windows" {
		comspec := os.Getenv("ComSpec")
		comspecLower := strings.ToLower(comspec)
		switch {
		case strings.Contains(comspecLower, "powershell"):
			shell = "PowerShell"
		case strings.Contains(comspecLower, "cmd"):
			shell = "CMD"
		default:
			shell = filepath.Base(comspec)
		}
	} else {
		shellEnv := os.Getenv("SHELL")
		if shellEnv != "" {
			shell = filepath.Base(shellEnv) // e.g., bash, zsh, etc.
		} else {
			shell = "unknown"
		}
	}

	return fmt.Sprintf("%s/%s", goos, shell)

}
