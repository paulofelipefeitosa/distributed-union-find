package master

import (
	"io/ioutil"
	"log"
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/paulofelipefeitosa/distributed-union-find/utils"
	"github.com/paulofelipefeitosa/distributed-union-find/config"
	"github.com/paulofelipefeitosa/distributed-union-find/protocol"
)

type ServerMaster struct {
	ID int
	Neighbors []string
	GrabbedNeighborhood []int //index
	DomainSize int
	URL config.URL
	EdgeMasters utils.ConcurrentSet
}

func (server *ServerMaster) Build(appConfig config.AppConfig) {
	server.ID = 0
	server.Neighbors = appConfig.Neighborhood()
	server.DomainSize = 1
	server.URL = appConfig.URL()
	server.EdgeMasters = utils.ConcurrentSet{}

	// TODO: Try to Grab
}

func (server *ServerMaster) grab(url string) bool {
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(protocol.GrabRequest{Initiator: server.URL})

	resp, err := http.Post(url, "application/json", reqBody)
	if err != nil {
		log.Printf("Cannot send grab request to %s due to %+v\n", url, err)
		return false
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(resp.Body)
		grabStats := protocol.GrabResponse{}
		err = decoder.Decode(&grabStats)
		if err != nil {
			log.Printf("Cannot decode response from grab request to %s due to %+v\n", url, err)
			return false
		}

		//TODO: response ok, if grabs > 0 then mark index as containing and increment domain size. Anyway merge Edge masters.
	} else {
		log.Printf("Send grab request to %s received %d status code\n", url, resp.StatusCode)
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Cannot read response from grab request to %s due to %+v\n", url, err)
		} else {
			log.Printf("Send grab request to %s received as response body: %s\n", url, respBody)
		}
		return false
	}
}

func (server *ServerMaster) GetURL() config.URL {
	return server.URL
}

func (server *ServerMaster) InfoFuncHandler() http.HandlerFunc {
	type InfoServerMaster struct {
		ID int
		Neighbors []string
		URL config.URL
		EdgeMasters []interface{}
	}
	return func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		body := toBytes(InfoServerMaster{ID:server.ID, Neighbors:server.Neighbors, URL:server.URL, EdgeMasters: server.EdgeMasters.Keys()})
		w.Write(body)
	}
}

func (server *ServerMaster) GrabFuncHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		grabberServer := protocol.GrabRequest{}
		err := decoder.Decode(&grabberServer)
		if err != nil {
			log.Printf("Cannot decode request body %+v\n", err)
			return
		}
		log.Printf("Adding %+v to Edge Masters and returning %+v public address\n", grabberServer, server.URL)
		server.EdgeMasters.Insert(grabberServer)

		w.WriteHeader(http.StatusOK)
		body := toBytes(protocol.GrabResponse{Grabs: 0, EdgeInitiators: []config.URL{server.URL}})
		w.Write(body)
	}
}

func toBytes(v interface{}) []byte {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(v)
	return body.Bytes()
}
