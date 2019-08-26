package protocol

import "github.com/paulofelipefeitosa/distributed-union-find/config"

type GrabResponse struct {
	Grabs int
	EdgeInitiators []config.URL
}

type GrabRequest struct {
	Initiator config.URL
}
