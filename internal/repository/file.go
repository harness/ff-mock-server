package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/drone/ff-mock-server/pkg/api"
)

const source = "./tests"

type memoryStore struct {
	configs  []api.FeatureConfig
	segments []api.Segment
}

type fileStruct struct {
	flag     api.FeatureConfig `json:"flag"`
	segments []api.Segment     `json:"segments"`
}

func init() {
	loadFiles()
}

func loadFiles() {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		log.Print(err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

	}
}

func loadFile(filename string) fileStruct {
	content, err := ioutil.ReadFile(filepath.Join(source, filename))
	if err != nil {
		log.Print(err)
	}

	result := fileStruct{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		log.Print(err)
	}
	return result
}
