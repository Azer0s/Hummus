package main

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/project"
	"github.com/Azer0s/Hummus/runner"
	log "github.com/sirupsen/logrus"
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
			log.SetLevel(log.TraceLevel)

			currentDir, err := os.Getwd()

			if err != nil {
				log.Fatal(err.Error())
			}

			project.BuildProject(currentDir)
		} else if os.Args[1] == "run" {
			currentDir, err := os.Getwd()

			if err != nil {
				panic(err)
			}

			project.RunProject(currentDir)
		} else {
			runner.RunFile(os.Args[1])
		}
	} else {
		runner.RunRepl()
	}
}
