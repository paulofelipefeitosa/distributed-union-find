package master

import (
	"bytes"
	"encoding/json"
	"github.com/paulofelipefeitosa/distributed-union-find/config"
	"github.com/paulofelipefeitosa/distributed-union-find/protocol"
	"github.com/paulofelipefeitosa/distributed-union-find/utils"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ServerMaster struct {
	ID int
	Neighbors []string
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

	neighbors := len(server.Neighbors)
	gbStatsCh := make(chan protocol.GrabResponse, neighbors)
	for _, element := range server.Neighbors {
		go server.Grab(element, gbStatsCh)
	}
	go server.MergeGrabs(gbStatsCh, neighbors)
}

func (server *ServerMaster) MergeGrabs(ch chan protocol.GrabResponse, neighbors int) {
	for i := 0; i < neighbors;i++ {
		grabStats := <- ch
		server.DomainSize += grabStats.Grabs
		for _, element := range grabStats.EdgeInitiators {
			server.EdgeMasters.Insert(element)
		}
	}
}

func (server *ServerMaster) Grab(url string, ch chan protocol.GrabResponse) {
	reqBody := new(bytes.Buffer)
	json.NewEncoder(reqBody).Encode(protocol.GrabRequest{Initiator: server.URL})

	for {
		resp, err := http.Post(url, "application/json", reqBody)
		if err != nil {
			log.Printf("Cannot send Grab request to %s due to %+v\n", url, err)
			resp.Body.Close()
			break
		}

		if resp.StatusCode == http.StatusOK {
			decoder := json.NewDecoder(resp.Body)
			grabStats := protocol.GrabResponse{}
			err = decoder.Decode(&grabStats)
			resp.Body.Close()
			if err != nil {
				log.Printf("Cannot decode response from Grab request to %s due to %+v\n", url, err)
				break
			}

			ch <- grabStats
			break
		} else {
			log.Printf("Send Grab request to %s received %d status code\n", url, resp.StatusCode)
			respBody, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("Cannot read response from Grab request to %s due to %+v\n", url, err)
			} else {
				log.Printf("Send Grab request to %s received as response body: %s\n", url, respBody)
			}
		}
		log.Printf("Sleeping 5 secs to retry %s Grabbing", url)
		time.Sleep(time.Second * 5)
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
		server.EdgeMasters.Insert(grabberServer.Initiator)

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
