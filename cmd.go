package main

import (
	"os"
	"os/exec"
)

func execute(command string) error {
	cmd := exec.Command("sh", []string{"-c", command}...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
