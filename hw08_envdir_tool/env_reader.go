package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := make(Environment, len(files))

	for _, file := range files {
		if !isValid(file) {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		value, err := getValue(dir, file.Name())
		if err != nil {
			continue
		}

		envs[file.Name()] = EnvValue{
			Value:      clean(value),
			NeedRemove: info.Size() == 0,
		}
	}

	return envs, nil
}

func isValid(file os.DirEntry) bool {
	if file.IsDir() {
		return false
	}

	if strings.Contains(file.Name(), "=") {
		return false
	}

	return true
}

func clean(s string) string {
	name := strings.TrimRight(s, " \t")
	return strings.ReplaceAll(name, "\x00", "\n")
}

func getValue(dir string, fileName string) (string, error) {
	f, err := os.Open(dir + "/" + fileName)
	if err != nil {
		return "", err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return "", err
		}
		return "", nil
	}

	return scanner.Text(), err
}
