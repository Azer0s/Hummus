package project

import (
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/Azer0s/Hummus/runner"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

//RunProject runs a Hummus project and builds it if it wasn't built before
func RunProject(projectDir string) {
	settings := readSettings(path.Join(projectDir, "project.json"))

	if _, err := os.Stat(path.Join(projectDir, settings.Output)); os.IsNotExist(err) {
		log.Warn("Project was not built! Building...")
		BuildProject(projectDir)
	}

	interpreter.LibBasePath = path.Join(projectDir, "lib/")

	runner.RunFile(path.Join(projectDir, settings.Output, settings.Entry))
}
