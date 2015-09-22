package common

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// WebServiceHost is a container for multiple Service stubs to live in
type WebServiceHost struct {
	martini *martini.Martini
	stubs   []ServiceStub
}

// Start starts up the webserver and begins serving request to the stubs
func (wsh *WebServiceHost) Start(URI string) {
	var m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(render.Renderer())

	for _, element := range wsh.stubs {
		m.Action(element.Route().Handle)
	}
	wsh.martini = m
	wsh.martini.RunOnAddr(URI)
}

// Stop terminates the webserver
func (wsh *WebServiceHost) Stop() {

}

// Use registers a stub with the service host
func (wsh *WebServiceHost) Use(stub ServiceStub) {
	wsh.stubs = append(wsh.stubs, stub)
}

// NewServiceHost creates a new initialized service host
func NewServiceHost() WebServiceHost {
	return WebServiceHost{}
}
