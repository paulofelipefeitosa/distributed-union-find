package config

import (
	"encoding/json"
	"os"
)

const JsonExtension = ".json"

type JSONConfigurator struct {
	Initiator bool
	Neighbors []string
	MyIP string
	Port int
}

func (r JSONConfigurator) IsInitiator() bool {
	return r.Initiator
}

func (r JSONConfigurator) Neighborhood() []string {
	return r.Neighbors
}

func (r JSONConfigurator) URL() URL {
	return URL{IP: r.MyIP, Port: r.Port}
}

func (r JSONConfigurator) Read(confFilePath string) (AppConfig, error) {
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