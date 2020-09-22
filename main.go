package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/project"
	"github.com/Azer0s/Hummus/runner"
	"log"
	"os"
	"path"
)

func init() {
	p, err := os.Executable()

	if err != nil {
		panic(err)
	}

	interpreter.BasePath = path.Dir(p)
	interpreter.LibBasePath = path.Join(interpreter.BasePath, "lib")

	err = os.Mkdir(interpreter.LibBasePath, os.ModePerm)

	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err.Error())
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "init" {
			project.InitProject(os.Args[2])
		} else if os.Args[1] == "build" {
			project.BuildProject()
		} else if os.Args[1] == "run" {
			project.RunProject()
		} else {
			runner.RunFile(os.Args[1])
		}
	} else {
		runner.RunRepl()
	}
}
