package apisrv

import "github.com/enzian/go-msf/common"

// RouteCache keeps a list of combination between a key (version and service prefix) and the URI to which the request will be dispatched
type RouteCache map[string]([]string)

// Projection is the handler for a specific Event
type Projection func(RouteCache, common.Event) (RouteCache, error)

// A MessageAdapter produces messages for the projection queue.
type MessageAdapter interface {
	Start(sendEventChan chan<- common.Event)
	Stop()
}
