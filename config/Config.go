package config

import (
	"errors"
	"fmt"
	"strings"
)

type AppConfig interface {
	IsInitiator() bool
	Neighborhood() []string
	URL() URL
}

type Configurator interface {
	Read(confFilePath string) (AppConfig, error)
}

type DefaultConfigurator struct {

}

func (r DefaultConfigurator) Read(confFilePath string) (AppConfig, error) {
	var config Configurator
	if strings.HasSuffix(confFilePath, JsonExtension) {
		config = JSONConfigurator{}
	} else {
		return nil, errors.New(fmt.Sprintf("Cannot identify config file type decoder [%s]", confFilePath))
	}
	return config.Read(confFilePath)
}

