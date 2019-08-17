package main

import (
	"fmt"
	"log"
	"github.com/paulofelipefeitosa/distributed-union-find/config"
)

func main() {
	conf, err := DefaultConfigurator{}.read("appconfig.json")
	if err != nil {
		log.Fatal(err)
	}
	if conf.isInitiator() {
		fmt.Printf("is initiator")
	} else {
		fmt.Printf("is not initiator")
	}
}