package codeexec

import (
	"errors"
	"os/exec"
	"strings"

	"open-sandbox/internal/shell"
)

func Exec(runtime string, args []string, workingDir string) (shell.Result, error) {
	command := strings.ToLower(strings.TrimSpace(runtime))
	if command == "" {
		return shell.Result{}, errors.New("runtime is required")
	}

	var binary string
	switch command {
	case "python":
		binary = "python"
	case "node":
		binary = "node"
	default:
		return shell.Result{}, errors.New("unsupported runtime")
	}

	if _, err := exec.LookPath(binary); err != nil {
		return shell.Result{}, err
	}

	return shell.Exec(binary, args, workingDir)
}
