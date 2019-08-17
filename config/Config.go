package config

import (
	"errors"
	"fmt"
	"strings"
)

type AppConfig interface {
	isInitiator() bool
	neighborhood() []string
}

type Configurator interface {
	read(confFilePath string) (AppConfig, error)
}

type DefaultConfigurator struct {

}

func (r DefaultConfigurator) read(confFilePath string) (AppConfig, error) {
	var config Configurator
	if strings.HasSuffix(confFilePath, JsonExtension) {
		config = JSONConfigurator{}
	} else {
		return nil, errors.New(fmt.Sprintf("Cannot identify config file type decoder [%s]", confFilePath))
	}
	return config.read(confFilePath)
}
