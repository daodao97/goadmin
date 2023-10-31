package util

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func NewCommand(dir, cmd string, args ...string) (*exec.Cmd, *bytes.Buffer) {
	buf := &bytes.Buffer{}

	toolCmd := &exec.Cmd{
		Path:   cmd,
		Args:   append([]string{cmd}, args...),
		Dir:    dir,
		Stdin:  os.Stdin,
		Stdout: buf,
		Stderr: buf,
		Env:    os.Environ(),
	}
	if filepath.Base(cmd) == cmd {
		if lp, err := exec.LookPath(cmd); err == nil {
			toolCmd.Path = lp
		}
	}
	return toolCmd, buf
}

func Command(dir, cmd string, args ...string) (output string, err error) {
	toolCmd, buf := NewCommand(dir, cmd, args...)
	err = toolCmd.Run()
	output = buf.String()
	return
}

func StdCommand(ctx context.Context, dir, cmd string, args ...string) (err error) {
	toolCmd := &exec.Cmd{
		Path:   cmd,
		Args:   append([]string{cmd}, args...),
		Dir:    dir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}
	if filepath.Base(cmd) == cmd {
		if lp, err := exec.LookPath(cmd); err == nil {
			toolCmd.Path = lp
		}
	}
	if err = toolCmd.Start(); err != nil {
		return
	}

	stop := make(chan error, 1)
	go func() {
		stop <- toolCmd.Wait()
		close(stop)
	}()

	select {
	case <-ctx.Done():
		toolCmd.Process.Signal(syscall.SIGINT)
		err = <-stop
	case err = <-stop:
	}
	return
}
