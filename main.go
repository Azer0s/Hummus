package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/runner"
	"os"
	"path"
)

func init() {
	p, err := os.Executable()

	if err != nil {
		panic(err)
	}

	interpreter.BasePath = path.Dir(p)
}

func main() {
	if len(os.Args) > 1 {
		runner.RunFile(os.Args[1])
	} else {
		runner.RunRepl()
	}
}
