package project

import (
	"encoding/json"
	"io/ioutil"
)

type projectJson struct {
	Name     string   `json:"name"`
	Packages []string `json:"packages"`
	Output   string   `json:"output"`
	Exclude  []string `json:"exclude"`
	Native   []string `json:"native"`
	Entry    string   `json:"entry"`
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
