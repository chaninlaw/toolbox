package utils

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
)

func NewError(e interface{}) error {
	if e != nil {
		_, fn, line, _ := runtime.Caller(1)
		return fmt.Errorf("[error] %s:%d %v", relative(fn), line, e)
	}
	return nil
}

// ExecCommand executes a command in the given relative directory ("." for current directory)
func ExecCommandInDir(dir string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	var stderrBuf []byte
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}
	stderrBuf, _ = io.ReadAll(stderrPipe)
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("command failed: %w\nstderr: %s", err, string(stderrBuf))
	}
	return nil
}
