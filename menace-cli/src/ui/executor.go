package ui

import (
	"os/exec"
	"runtime"
)

// Runs a shell command
func runShellCommand(command string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	output, err := cmd.CombinedOutput()
	return string(output), err
}
