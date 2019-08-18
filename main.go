package main

import (
	"github.com/paulofelipefeitosa/distributed-union-find/config"
	"log"
)

func main() {
	conf, err := config.DefaultConfigurator{}.Read("appconfig.json")
	if err != nil {
		log.Fatal(err)
	}
	Serve(conf)
}