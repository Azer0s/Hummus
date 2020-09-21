package project

import (
	"encoding/json"
	"github.com/Azer0s/Hummus/interpreter"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var builtPackages = make([]string, 0)
var basePath string

func init() {
	p, err := os.Executable()

	if err != nil {
		panic(err)
	}

	basePath = path.Join(path.Dir(p), "..")
}

//BuildProject builds a Hummus project
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

	nativePackages := make([]string, 0)
	for _, n := range settings.Native {
		b, _ := json.Marshal(n)
		nativePackages = append(nativePackages, string(b))
	}

	log.Tracef("Native files: %s", "["+strings.Join(nativePackages, ", ")+"]")

	packages := make([]string, 0)
	for _, j := range settings.Packages {
		b, _ := json.Marshal(j)
		packages = append(packages, string(b))
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

func buildNativeLibs(currentDir string, nativeLibs []nativePackage, outputFolder string) {
	if !(len(nativeLibs) > 0) {
		return
	}

	log.Info("Building native files...")

	for _, lib := range nativeLibs {
		log.Debugf("Building native file %s...", lib.Out)

		tmpFolder := path.Join(basePath, strings.ReplaceAll(uuid.New().String(), "-", ""))

		log.Tracef("Creating temporary folder %s in Hummus root...", tmpFolder)

		err := os.Mkdir(tmpFolder, os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}

		//goland:noinspection ALL
		defer func() {
			log.Tracef("Removing temporary folder %s...", tmpFolder)
			err = os.RemoveAll(tmpFolder)

			if err != nil {
				log.Fatal(err.Error())
			}
		}()

		targetFiles := make([]string, 0)

		for _, file := range lib.Files {
			target := path.Join(tmpFolder, file)

			log.Tracef("Copying %s", file)

			err = os.MkdirAll(path.Dir(target), os.ModePerm)
			if err != nil {
				log.Fatal(err.Error())
			}

			err = copyFile(path.Join(currentDir, file), target)
			if err != nil {
				log.Fatal(err.Error())
			}

			targetFiles = append(targetFiles, path.Join(path.Base(tmpFolder), file))
		}

		//HACK: Due do delve, his does not work in debug
		cmd := exec.Command("go", "build", "-buildmode=plugin" /*, "-gcflags='all=-N -l'"*/, "-o", path.Join(path.Base(tmpFolder), lib.Out))
		cmd.Args = append(cmd.Args, targetFiles...)

		cmd.Dir = basePath

		log.Tracef("Running \"%s\"...", cmd.String())

		buff, err := cmd.CombinedOutput()

		if len(buff) > 0 {
			log.Tracef(interpreter.ReplaceEnd(strings.ReplaceAll(string(buff), "\n", "; "), "; ", "", -1))
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		outputPath := path.Join(outputFolder, lib.Out)
		err = os.MkdirAll(path.Dir(outputPath), os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}

		object := path.Join(tmpFolder, lib.Out)

		log.Tracef("Copying shared object (%s) to output path (%s)...", object, outputPath)

		err = copyFile(object, outputPath)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func copyFiles(currentDir string, nativeLibs []nativePackage, excludedFiles []string, outputFolder string) {
	log.Info("Copying files to output directory...")

	absoluteNativeLibs := make([]string, 0)
	for _, lib := range nativeLibs {
		for _, file := range lib.Files {
			absoluteNativeLibs = append(absoluteNativeLibs, path.Join(currentDir, file))
		}
	}

	absoluteExcludedFiles := make([]string, 0)
	for _, file := range excludedFiles {
		absoluteExcludedFiles = append(absoluteExcludedFiles, path.Join(currentDir, file))
	}

	absoluteBuiltNativeLibs := make([]string, 0)
	for _, lib := range nativeLibs {
		absoluteBuiltNativeLibs = append(absoluteBuiltNativeLibs, path.Join(currentDir, outputFolder, strings.ReplaceAll(lib.Out, ".go", ".so")))
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

			if contains(absoluteBuiltNativeLibs, filePath) {
				return nil
			}

			relLibPath, err := filepath.Rel(path.Join(currentDir, "lib"), filePath)
			relOutPath, err := filepath.Rel(path.Join(currentDir, outputFolder), filePath)

			if filePath == path.Join(currentDir, outputFolder) || filePath == path.Join(currentDir, "lib") ||
				!strings.Contains(relLibPath, "..") || !strings.Contains(relOutPath, "..") {
				return nil
			}

			if contains(absoluteNativeLibs, filePath) || contains(absoluteExcludedFiles, filePath) ||
				filePath == path.Join(currentDir, "project.json") {
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
					if os.IsExist(err) {
						return nil
					}
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

		log.Debugf("Pulling package %s@%s...", s.Repo, s.At)

		tmpFolder := strings.ReplaceAll(uuid.New().String(), "-", "")

		cmd := exec.Command("git", "clone", repoUrl, path.Join(libFolder, tmpFolder))
		buff, err := cmd.CombinedOutput()

		if len(buff) > 0 {
			log.Tracef(interpreter.ReplaceEnd(strings.ReplaceAll(string(buff), "\n", "; "), "; ", "", -1))
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		cmd = exec.Command("git", "checkout", s.At)

		cmd.Dir = path.Join(libFolder, tmpFolder)

		buff, err = cmd.CombinedOutput()

		if len(buff) > 0 {
			log.Tracef(interpreter.ReplaceEnd(strings.ReplaceAll(string(buff), "\n", "; "), "; ", "", -1))
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		libSettings := readSettings(path.Join(libFolder, tmpFolder, "project.json"))

		name := libSettings.Name

		if s.At != "master" {
			name += "@" + strings.ReplaceAll(s.At, ".", "_")
		}

		tmpPath := path.Join(libFolder, tmpFolder)
		err = os.Rename(tmpPath, path.Join(libFolder, name))

		if err != nil {
			if os.IsExist(err) {
				log.Warnf("%s has already been pulled!", name)
				log.Debugf("Deleting %s...", tmpPath)

				_ = os.RemoveAll(tmpPath)

				builtPackages = append(builtPackages, repoUrl+"@"+s.At)
				continue
			} else {
				log.Fatal(err.Error())
			}
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

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
