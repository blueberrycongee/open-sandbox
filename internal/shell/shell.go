package shell

import (
	"bytes"
	"os/exec"
)

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Exec(command string, args []string, workingDir string) (Result, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workingDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return Result{
				Stdout:   stdout.String(),
				Stderr:   stderr.String(),
				ExitCode: exitErr.ExitCode(),
			}, nil
		}
		return Result{}, err
	}

	return Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
	}, nil
}
