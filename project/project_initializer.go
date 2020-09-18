package project

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const readme = `# {project_name}

This is the README of {project_name}

## Getting started

` + "```bash" + `
hummus build
hummus run
` + "```" + `

## License

This is the place for your license information
`

const gitignore = `### Intellij ###
# Covers JetBrains IDEs: IntelliJ, RubyMine, PhpStorm, AppCode, PyCharm, CLion, Android Studio and WebStorm
# Reference: https://intellij-support.jetbrains.com/hc/en-us/articles/206544839

# User-specific stuff
.idea/**/workspace.xml
.idea/**/tasks.xml
.idea/**/usage.statistics.xml
.idea/**/dictionaries
.idea/**/shelf

# Generated files
.idea/**/contentModel.xml

# Sensitive or high-churn files
.idea/**/dataSources/
.idea/**/dataSources.ids
.idea/**/dataSources.local.xml
.idea/**/sqlDataSources.xml
.idea/**/dynamic.xml
.idea/**/uiDesigner.xml
.idea/**/dbnavigator.xml

# Gradle
.idea/**/gradle.xml
.idea/**/libraries

# CMake
cmake-build-*/

# Mongo Explorer plugin
.idea/**/mongoSettings.xml

# File-based project format
*.iws

# IntelliJ
out/

# mpeltonen/sbt-idea plugin
.idea_modules/

# JIRA plugin
atlassian-ide-plugin.xml

# Cursive Clojure plugin
.idea/replstate.xml

# Crashlytics plugin (for Android Studio and IntelliJ)
com_crashlytics_export_strings.xml
crashlytics.properties
crashlytics-build.properties
fabric.properties

# Editor-based Rest Client
.idea/httpRequests

# Android studio 3.1+ serialized cache file
.idea/caches/build_file_checksums.ser

### Intellij Patch ###
# Comment Reason: https://github.com/joeblau/gitignore.io/issues/186#issuecomment-215987721

# *.iml
# modules.xml
# .idea/misc.xml
# *.ipr

# Sonarlint plugin
.idea/**/sonarlint/

# SonarQube Plugin
.idea/**/sonarIssues.xml

# Markdown Navigator plugin
.idea/**/markdown-navigator.xml
.idea/**/markdown-navigator/

### macOS ###
# General
.DS_Store
.AppleDouble
.LSOverride

### Hummus ###

# Build output
bin/
`

const main = `(use :<base>)

(out "Hello, World!")
`

func gitAdd(filename, projectDir string) {
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = projectDir
	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}

//InitProject initializes a Hummus project
func InitProject(projectName string) {
	currentDir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	projectName = strings.ReplaceAll(projectName, " ", "_")
	projectDir, err := filepath.Abs(path.Join(currentDir, projectName))

	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		panic(projectDir + " already exists!")
	}

	err = os.Mkdir(projectDir, os.ModePerm)

	if err != nil {
		panic(err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = projectDir
	err = cmd.Run()

	if err != nil {
		panic(err)
	}

	settings := projectJson{
		Name:     projectName,
		Packages: make([]packageJson, 0),
		Output:   "bin",
		Exclude: []string{
			".gitignore",
			"README.md",
		},
		Native: make([]nativePackage, 0),
		Entry:  "main.hummus",
	}

	b, err := json.MarshalIndent(settings, "", "    ")

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(projectDir, "project.json"), b, os.ModePerm)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(projectDir, "README.md"), []byte(strings.ReplaceAll(readme, "{project_name}", projectName)), os.ModePerm)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(projectDir, ".gitignore"), []byte(gitignore), os.ModePerm)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(projectDir, "main.hummus"), []byte(main), os.ModePerm)

	if err != nil {
		panic(err)
	}

	gitAdd("project.json", projectDir)
	gitAdd("README.md", projectDir)
	gitAdd(".gitignore", projectDir)
	gitAdd("main.hummus", projectDir)

	cmd = exec.Command("git", "commit", "-m", "\"Initial commit\"")
	cmd.Dir = projectDir
	err = cmd.Run()

	if err != nil {
		panic(err)
	}
}
