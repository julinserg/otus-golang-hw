package main

import (
	"bufio"
	"bytes"
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

func readFirstLineInFile(pathToFile string) (*EnvValue, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileStat, err := os.Stat(pathToFile)
	if err != nil {
		return nil, err
	}
	if fileStat.Size() == 0 {
		return &EnvValue{"", true}, err
	}

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("ReadDir error scanner.Scan() - scanner not read line")
	}
	firstStrFromFile := scanner.Text()
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &EnvValue{firstStrFromFile, false}, nil
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
		if strings.Contains(v.Name(), "=") {
			continue
		}
		pathToFile := filepath.Join(dir, v.Name())
		valEnvRaw, err := readFirstLineInFile(pathToFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if valEnvRaw.NeedRemove {
			env[v.Name()] = *valEnvRaw
			continue
		}
		dat := []byte(valEnvRaw.Value)
		fromRep := []byte{0}
		toRep := []byte("\n")
		dat = bytes.ReplaceAll(dat, fromRep, toRep)
		strDat := string(dat)
		valEnv := strings.TrimRight(strDat, " 	")

		env[v.Name()] = EnvValue{valEnv, false}
	}
	return env, nil
}
