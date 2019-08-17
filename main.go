package main

import (
	"fmt"
	"github.com/paulofelipefeitosa/distributed-union-find/config"
	"log"
)

func main() {
	conf, err := config.DefaultConfigurator{}.Read("appconfig.json")
	if err != nil {
		log.Fatal(err)
	}
	if conf.IsInitiator() {
		fmt.Printf("is initiator\n")
	} else {
		fmt.Printf("is not initiator\n")
	}
}