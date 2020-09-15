package project

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func BuildProject() {
	log.SetLevel(log.TraceLevel)

	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatal(err.Error())
	}

	settingsPath := path.Join(currentDir, "project.json")
	log.Debugf("Reading project settings '%s'...", settingsPath)
	settings := readSettings(settingsPath)

	log.Infof("Building project %s", settings.Name)

	log.Tracef("Project name: %s", settings.Name)
	log.Tracef("Entry point: %s", settings.Entry)
	log.Tracef("Output path: %s", settings.Output)
	log.Tracef("Excluded files: %s", "["+strings.Join(settings.Exclude, ", ")+"]")
	log.Tracef("Native files: %s", "["+strings.Join(settings.Native, ", ")+"]")

	packages := make([]string, 0)
	for _, json := range settings.Packages {
		packages = append(packages, json.String())
	}

	log.Tracef("Packages: %s", "["+strings.Join(packages, ", ")+"]")

	createOutputFolder(path.Join(currentDir, settings.Output))
	createLibFolder(path.Join(currentDir, "lib/"))
	buildNativeLibs(currentDir, settings.Native, settings.Output)
	copyFiles(settings.Native, settings.Exclude, settings.Output)
	pullPackages(path.Join(currentDir, "lib/"), settings.Packages)

	log.Info("Project built successfully!")
}

func createOutputFolder(folder string) {
	log.Debugf("Creating output folder '%s'...", folder)

	if _, err := os.Stat(folder); !os.IsNotExist(err) {
		log.Warn("Output folder already exists! Wiping...")

		err = os.RemoveAll(folder)

		if err != nil {
			log.Fatal(err.Error())
		}
	}

	err := os.Mkdir(folder, os.ModePerm)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func createLibFolder(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		log.Warn("Lib folder does not exist! Creating...")
		log.Debugf("Creating lib folder '%s'...", folder)

		err := os.Mkdir(folder, os.ModePerm)

		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func buildNativeLibs(currentDir string, nativeLibs []string, outputFolder string) {
	if len(nativeLibs) > 0 {
		log.Info("Building native files...")
	}

	for _, lib := range nativeLibs {
		log.Debugf("Building native file %s...", lib)
		cmd := exec.Command("go", "build", "-buildmode=plugin", "-gcflags='all=-N -l'", "-o", path.Join(outputFolder, strings.Replace(lib, ".go", ".so", 1)), path.Join(currentDir, lib))

		err := cmd.Start()

		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func copyFiles(nativeLibs []string, excludedFiles []string, outputFolder string) {
	log.Info("Copying files to output directory...")

	err := filepath.Walk(".",
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.Index(filePath, ".") == 0 {
				log.Tracef("Skipping hidden file %s", filePath)
				return nil
			}

			if contains(nativeLibs, filePath) || contains(excludedFiles, filePath) || filePath == outputFolder || filePath == "project.json" || filePath == "lib" {
				if !info.IsDir() {
					log.Tracef("Skipping file %s", filePath)
				} else {
					log.Tracef("Skipping directory %s/", filePath)
				}

				return nil
			}

			log.Debugf("Copying %s", filePath)

			input, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatal(err.Error())
			}

			err = ioutil.WriteFile(path.Join(outputFolder, filePath), input, 0644)
			if err != nil {
				log.Fatal(err.Error())
			}

			return nil
		})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func pullPackages(libFolder string, packages []packageJson) {
	log.Info("Pulling packages...")

	for _, s := range packages {
		log.Debugf("Pulling package %s...", s)

		cmd := exec.Command("git", "clone", "https://"+s.Repo)
		cmd.Dir = libFolder
		err := cmd.Start()

		if err != nil {
			log.Fatal(err.Error())
		}

		if s.At == "master" {
			continue
		}

		cmd = exec.Command("git", "checkout", s.At)

		repo := strings.Split(s.Repo, "/")
		cmd.Dir = path.Join(libFolder, repo[len(repo)-1])

		err = cmd.Start()

		if err != nil {
			log.Fatal(err.Error())
		}

		err = os.Rename(path.Join(libFolder, repo[len(repo)-1]), path.Join(libFolder, repo[len(repo)-1]+"@"+s.At))

		if err != nil {
			log.Fatal(err.Error())
		}
	}

	//TODO: Build packages recursively
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
