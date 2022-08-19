package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	f, err := os.Open(dir)
	if err != nil {
		fmt.Println("ReadDir error os.Open dir -" + err.Error())
		return nil, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println("ReadDir error f.Readdir -" + err.Error())
		return nil, err
	}

	env := make(Environment)
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		pathToFile := filepath.Join(dir, v.Name())
		file, err := os.Open(pathToFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer file.Close()

		fileStat, err := os.Stat(pathToFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if fileStat.Size() == 0 {
			env[v.Name()] = EnvValue{"", true}
			continue
		}

		scanner := bufio.NewScanner(file)
		if !scanner.Scan() {
			fmt.Println("ReadDir error scanner.Scan() - scanner not read line")
			continue
		}
		valEnv := scanner.Text()
		strings.TrimRight(valEnv, " 	")
		strings.ReplaceAll(valEnv, `0x00`, `\n`)
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			continue
		}

		env[v.Name()] = EnvValue{valEnv, false}
	}
	return env, nil
}
