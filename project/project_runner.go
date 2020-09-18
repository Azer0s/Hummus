package project

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/runner"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

//RunProject runs a Hummus project and builds it if it wasn't built before
func RunProject() {
	currentDir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	settings := readSettings(path.Join(currentDir, "project.json"))

	if _, err := os.Stat(path.Join(currentDir, settings.Output)); os.IsNotExist(err) {
		log.Warn("Project was not built! Building...")
		BuildProject()
	}

	interpreter.LibBasePath = path.Join(currentDir, "lib/")

	runner.RunFile(path.Join(currentDir, settings.Output, settings.Entry))
}
