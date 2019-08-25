package master

import (
	"bytes"
	"encoding/json"
	"github.com/paulofelipefeitosa/distributed-union-find/config"
	"github.com/paulofelipefeitosa/distributed-union-find/utils"
	"log"
	"net/http"
)

type ServerMaster struct {
	ID int
	Neighbors []string
	URL config.URL
	EdgeMasters utils.ConcurrentSet
}

type PublicMasterServer struct {
	URL config.URL
}

func (server *ServerMaster) Build(appConfig config.AppConfig) {
	server.ID = 0
	server.Neighbors = appConfig.Neighborhood()
	server.URL = appConfig.URL()
	server.EdgeMasters = utils.ConcurrentSet{}
}

func (server *ServerMaster) GetURL() config.URL {
	return server.URL
}

func (server *ServerMaster) InfoFuncHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		type InfoServerMaster struct {
			ID int
			Neighbors []string
			URL config.URL
			EdgeMasters []interface{}
		}

		body := toBytes(InfoServerMaster{ID:server.ID, Neighbors:server.Neighbors, URL:server.URL, EdgeMasters: server.EdgeMasters.Keys()})
		w.Write(body)
	}
}

func (server *ServerMaster) GrabFuncHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		grabberServer := PublicMasterServer{}
		err := decoder.Decode(&grabberServer)
		if err != nil {
			log.Printf("Cannot decode request body %+v\n", err)
			return
		}
		log.Printf("Adding %+v to Edge Masters and returning %+v public address\n", grabberServer, server.URL)
		server.EdgeMasters.Insert(grabberServer)

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
