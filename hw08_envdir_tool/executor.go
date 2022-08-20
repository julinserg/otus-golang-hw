package main

import (
	"errors"
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
	pathToCmd := cmd[0]
	if len(pathToCmd) == 0 {
		fmt.Println("name command is empty")
		return 1
	}
	execCmd := exec.Command(pathToCmd, cmd[1:]...)

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

	for key, value := range env {
		if value.NeedRemove {
			delete(envGlobalMap, key)
		}
	}

	for key, value := range env {
		if strings.Contains(key, "=") {
			continue
		}
		envGlobalMap[key] = value.Value
	}

	envGlobalResult := make([]string, 0, len(envGlobalMap))
	for key, value := range envGlobalMap {
		envGlobalResult = append(envGlobalResult, key+"="+value)
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
	if err := execCmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
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
