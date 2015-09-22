package apisrv

import "github.com/enzian/go-msf/api-server.back/services"

// RouteCache keeps a list of combination between a key (version and service prefix) and the URI to which the request will be dispatched
type RouteCache map[string]([]string)

// Event is a structre which represents a change to the route cache in some way.
type Event struct {
	Action string
	Data   map[string]string
}

// Projection is the handler for a specific Event
type Projection func(RouteCache, Event) (RouteCache, error)

// A MessageAdapter produces messages for the projection queue.
type MessageAdapter interface {
	Start(sendEventChan chan<- apisrv.Event)
	Stop()
}
