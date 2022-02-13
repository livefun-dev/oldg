package graph

import "github.com/lifevun-dev/oldg/graph/model"

type Observer struct {
	msgChan chan model.Command
}

var observers map[string]Observer = map[string]Observer{}
