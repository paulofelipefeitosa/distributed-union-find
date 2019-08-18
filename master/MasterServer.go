package master

import (
	"bytes"
	"encoding/json"
	"github.com/paulofelipefeitosa/distributed-union-find/config"
	"log"
	"net/http"
)

type ServerMaster struct {
	ID int
	Neighbors []string
	URL config.URL
	EdgeMasters []config.URL
}

type PublicMasterServer struct {
	URL config.URL
}

func (server *ServerMaster) Build(appConfig config.AppConfig) {
	(*server).ID = 0
	(*server).Neighbors = appConfig.Neighborhood()
	(*server).URL = appConfig.URL()
}

func (server ServerMaster) GetURL() config.URL {
	return server.URL
}

func (server ServerMaster) GetIDFuncHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.URL)
		w.WriteHeader(http.StatusOK)
		body := toBytes(server.ID)
		w.Write(body)
	}
}

func (server ServerMaster) GrabFuncHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		log.Printf("%+v\n", r.URL)

		// TODO: get MasterURL and add to edge masters initiators

		w.WriteHeader(http.StatusForbidden)
		body := toBytes(PublicMasterServer{URL: server.URL})
		w.Write(body)
	}
}

func toBytes(v interface{}) []byte {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(v)
	return body.Bytes()
}
