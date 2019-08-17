package config

import (
	"encoding/json"
	"os"
)

const JsonExtension = ".json"

type JSONConfigurator struct {
	Initiator bool
	Neighborhood []string
}

func (r JSONConfigurator) isInitiator() bool {
	return r.Initiator
}

func (r JSONConfigurator) neighborhood() []string {
	return r.Neighborhood
}

func (r JSONConfigurator) read(confFilePath string) (AppConfig, error) {
	file, err := os.Open(confFilePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	configuration := JSONConfigurator{}
	err = decoder.Decode(&configuration)
	return configuration, err
}