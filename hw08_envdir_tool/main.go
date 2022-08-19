package main

import (
	"log"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	retCode := RunCmd(os.Args[2:], env)
	os.Exit(retCode)
}
