package main

import (
	"os/exec"
	"strings"
)

// Handles the command passed in by -C
func prepare_cmd(input string) *exec.Cmd {
	chopped := strings.Split(input, " ")

	if len(chopped) < 1 {
		exit(1)
	}

	return exec.Command(chopped[0], chopped[1:]...)
}

// Renew the given command - that is, produce a new command using the same Path
// and Args from the given command.
func renew_cmd(old_cmd *exec.Cmd) *exec.Cmd {
	return &exec.Cmd{Path: old_cmd.Path, Args: old_cmd.Args}
}
