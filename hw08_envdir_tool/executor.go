package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		log.Fatal("name command not pass")
	}

	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Env = append(os.Environ(),
		"USER=petya",
		"CITY=SPb",
	)

	execCmd.Stdout = os.Stdout

	err := execCmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	if err := execCmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
				return status.ExitStatus()
			}
		} else {
			log.Fatalf("execCmd.Wait: %v", err)
		}
	}
	return
}
