package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		fmt.Println("name command not pass")
		return 1
	}

	execCmd := exec.Command(cmd[0], cmd[1:]...)

	envGlobal := os.Environ()
	envGlobalMap := make(map[string]string, len(envGlobal))
	for _, elem := range envGlobal {
		envStringSlice := strings.Split(elem, "=")
		if len(envStringSlice) == 1 {
			envGlobalMap[envStringSlice[0]] = ""
			continue
		}
		if len(envStringSlice) < 2 {
			continue
		}
		envGlobalMap[envStringSlice[0]] = envStringSlice[1]
	}

	for kEnv, vEnv := range env {
		if vEnv.NeedRemove {
			delete(envGlobalMap, kEnv)
		}
	}

	for kEnv, vEnv := range env {
		if strings.Contains(kEnv, "=") {
			continue
		}
		envGlobalMap[kEnv] = vEnv.Value
	}

	var envGlobalResult []string
	for kEnv, vEnv := range envGlobalMap {
		envGlobalResult = append(envGlobalResult, kEnv+"="+vEnv)
	}
	execCmd.Env = envGlobalResult

	execCmd.Stdout = os.Stdout
	execCmd.Stdin = os.Stdin
	execCmd.Stderr = os.Stderr

	err := execCmd.Start()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("Waiting for command to finish...")
	if err := execCmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("Exit Status: %d", status.ExitStatus())
				return status.ExitStatus()
			}
		} else {
			fmt.Printf("execCmd.Wait: %v", err)
			return 1
		}
	}
	return
}
