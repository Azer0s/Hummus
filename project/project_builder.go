package project

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var builtPackages = make([]string, 0)

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
	copyFiles(currentDir, settings.Native, settings.Exclude, settings.Output)
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

func copyFiles(currentDir string, nativeLibs []string, excludedFiles []string, outputFolder string) {
	log.Info("Copying files to output directory...")

	absoluteNativeLibs := make([]string, 0)
	for _, lib := range nativeLibs {
		absoluteNativeLibs = append(absoluteNativeLibs, path.Join(currentDir, lib))
	}

	absoluteExcludedFiles := make([]string, 0)
	for _, file := range excludedFiles {
		absoluteExcludedFiles = append(absoluteExcludedFiles, path.Join(currentDir, file))
	}

	err := filepath.Walk(currentDir,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			paths := strings.Split(filePath, "/")
			if anyHidden(paths) {
				if !info.IsDir() {
					log.Tracef("Skipping hidden file %s", filePath)
				} else {
					log.Tracef("Skipping hidden directory %s/", filePath)
				}

				return nil
			}

			if filePath == currentDir {
				return nil
			}

			if contains(absoluteNativeLibs, filePath) || contains(absoluteExcludedFiles, filePath) ||
				filePath == path.Join(currentDir, outputFolder) || filePath == path.Join(currentDir, "project.json") ||
				filePath == path.Join(currentDir, "lib") {
				if !info.IsDir() {
					log.Tracef("Skipping file %s", filePath)
				} else {
					log.Tracef("Skipping directory %s/", filePath)
				}

				return nil
			}

			log.Debugf("Copying %s", filePath)

			relPath, err := filepath.Rel(currentDir, filePath)

			if err != nil {
				log.Fatal(err.Error())
			}

			if info.IsDir() {
				err = os.Mkdir(path.Join(currentDir, outputFolder, relPath), os.ModePerm)

				if err != nil {
					log.Fatal(err.Error())
				}

				return nil
			}

			input, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatal(err.Error())
			}

			err = ioutil.WriteFile(path.Join(currentDir, outputFolder, relPath), input, os.ModePerm)
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
		repoUrl := "https://" + s.Repo

		if contains(builtPackages, repoUrl+"@"+s.At) {
			continue
		}

		log.Debugf("Pulling package %s...", s)

		tmpFolder := strings.ReplaceAll(uuid.New().String(), "-", "")

		cmd := exec.Command("git", "clone", repoUrl, path.Join(libFolder, tmpFolder))
		err := cmd.Start()

		if err != nil {
			log.Fatal(err.Error())
		}

		err = cmd.Wait()

		if err != nil {
			log.Fatal(err.Error())
		}

		cmd = exec.Command("git", "checkout", s.At)

		cmd.Dir = path.Join(libFolder, tmpFolder)

		err = cmd.Start()

		if err != nil {
			log.Fatal(err.Error())
		}

		err = cmd.Wait()

		if err != nil {
			log.Fatal(err.Error())
		}

		libSettings := readSettings(path.Join(libFolder, tmpFolder, "project.json"))

		name := libSettings.Name

		if s.At != "master" {
			name += "@" + s.At
		}

		err = os.Rename(path.Join(libFolder, tmpFolder), path.Join(libFolder, name))

		if err != nil {
			log.Fatal(err.Error())
		}

		builtPackages = append(builtPackages, repoUrl+"@"+s.At)

		log.Infof("Building library %s...", name)

		createOutputFolder(path.Join(libFolder, name, libSettings.Output))
		buildNativeLibs(path.Join(libFolder, name), libSettings.Native, libSettings.Output)
		copyFiles(path.Join(libFolder, name), libSettings.Native, libSettings.Exclude, libSettings.Output)
		pullPackages(libFolder, libSettings.Packages)
	}
}

func anyHidden(arr []string) bool {
	for _, elem := range arr {
		if len(elem) > 0 && elem[0] == '.' {
			return true
		}
	}

	return false
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
