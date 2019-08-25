package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paulofelipefeitosa/distributed-union-find/config"
	"github.com/paulofelipefeitosa/distributed-union-find/master"
	"log"
	"net/http"
)

type Component interface {
	Build(appConfig config.AppConfig)
	GetURL() config.URL

	GrabFuncHandler() http.HandlerFunc
	InfoFuncHandler() http.HandlerFunc
}

func Serve(conf config.AppConfig) {
	var component Component
	if conf.IsInitiator() {
		log.Print("Server is a Master App\n")
		component = &master.ServerMaster{}
		component.Build(conf)
	} else {
		// TODO:
		log.Print("Server is a Slave App\n")
	}

	router := mux.NewRouter()
	router.HandleFunc("/grab", component.GrabFuncHandler()).Methods(http.MethodPost)
	router.HandleFunc("/info", component.InfoFuncHandler()).Methods(http.MethodGet)

	serverURL := component.GetURL()
	log.Printf("Server URL %+v\n", serverURL)

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", serverURL.IP, serverURL.Port),
		Handler:        router,
	}

	log.Fatal(s.ListenAndServe())
}