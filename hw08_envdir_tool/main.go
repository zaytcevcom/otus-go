package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Invalid count of arguments")
		return
	}

	env, err := ReadDir(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	RunCmd(os.Args[2:], env)
}
