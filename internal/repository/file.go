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
	Flag     api.FeatureConfig `json:"flag"`
	Segments []api.Segment     `json:"segments"`
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
	fp := filepath.Clean(filepath.Join(source, filename))
	content, err := ioutil.ReadFile(fp)
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
