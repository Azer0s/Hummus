package main

import (
	"github.com/Azer0s/Hummus/runner"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		runner.RunFile(os.Args[1])
	} else {
		runner.RunRepl()
	}
}
