package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 2 {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	newEnv := getNewEnv(env)
	command.Env = make([]string, 0, len(newEnv))
	for key, value := range newEnv {
		command.Env = append(command.Env, key+"="+value)
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	command.Run()

	return command.ProcessState.ExitCode()
}

func getNewEnv(env Environment) map[string]string {
	envs := os.Environ()

	result := make(map[string]string, len(envs)+len(env))

	for _, envVar := range envs {
		envPair := strings.SplitN(envVar, "=", 2)
		result[envPair[0]] = envPair[1]
	}

	for key, value := range env {
		if value.NeedRemove {
			delete(result, key)
		} else {
			result[key] = value.Value
		}
	}

	return result
}
