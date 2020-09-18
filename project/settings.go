package project

import (
	"encoding/json"
	"io/ioutil"
)

type packageJson struct {
	Repo string `json:"repo"`
	At   string `json:"at"`
}

type nativePackage struct {
	Files []string `json:"files"`
	Out   string   `json:"out"`
}

type projectJson struct {
	Name     string          `json:"name"`
	Packages []packageJson   `json:"packages"`
	Output   string          `json:"output"`
	Exclude  []string        `json:"exclude"`
	Native   []nativePackage `json:"native"`
	Entry    string          `json:"entry"`
}

func readSettings(filename string) *projectJson {
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	settings := &projectJson{}

	err = json.Unmarshal(b, settings)

	if err != nil {
		panic(err)
	}

	return settings
}
