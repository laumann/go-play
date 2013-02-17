package main

import (
	"os/exec"
	"strings"
)

// Handles the command passed in by -C

func preparecmd(input string) *exec.Cmd {
	chopped := strings.Split(input, " ")

	if len(chopped) < 1 {
		exit(1)
	}

	return exec.Command(chopped[0], chopped[1:]...)
}

func renewcmd(old_cmd *exec.Cmd) *exec.Cmd {
	return &exec.Cmd{Path: old_cmd.Path, Args: old_cmd.Args}
}
