package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paulofelipefeitosa/distributed-union-find/config"
	"github.com/paulofelipefeitosa/distributed-union-find/master"
	"log"
	"net/http"
)

type Server interface {
	Build(appConfig config.AppConfig)
	GetURL() config.URL

	GrabFuncHandler() http.HandlerFunc
	GetIDFuncHandler() http.HandlerFunc
}

func Serve(conf config.AppConfig) {
	var server Server
	if conf.IsInitiator() {
		log.Print("Server is a Master App\n")
		server = &master.ServerMaster{}
		server.Build(conf)
	} else {
		log.Print("Server is a Slave App\n")
	}

	router := mux.NewRouter()
	router.HandleFunc("/grab/{MasterURL:.*}", server.GrabFuncHandler()).Methods(http.MethodGet)
	router.HandleFunc("/ID", server.GetIDFuncHandler()).Methods(http.MethodGet)

	serverURL := server.GetURL()
	log.Printf("Server URL %+v\n", serverURL)

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", serverURL.IP, serverURL.Port),
		Handler:        router,
	}

	log.Fatal(s.ListenAndServe())
}