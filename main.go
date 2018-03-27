package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func main() {
	var command string
	var args []string

	switch len(os.Args) {
	case 0, 1:
		exitErr("Usage: parent <command> [<args>...]")
	default:
		p, err := exec.LookPath(os.Args[1])
		if err != nil {
			exitErr(err)
		}
		command = p
		args = expandArgs(os.Args[2:])
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		exitErr(err)
	}

	go func() {
		for sig := range c {
			cmd.Process.Signal(sig)
		}
	}()

	if err := cmd.Wait(); err != nil {
		exitErr(err)
	}
}

func exitErr(errs ...interface{}) {
	fmt.Fprintln(os.Stderr, errs...)
	os.Exit(1)
}

// expandArgs leverages on shell and echo to expand
// possible args mainly env vars.
func expandArgs(args []string) []string {
	var expanded []string
	for _, arg := range args {
		e, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("echo %s", arg)).Output()
		// error is not expected.
		// in the rare case that this errors
		// the original arg is still used.
		if err == nil {
			arg = strings.TrimSpace(string(e))
		}
		expanded = append(expanded, arg)
	}
	return expanded
}
