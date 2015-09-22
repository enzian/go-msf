package common

import "github.com/go-martini/martini"

type ServiceStub interface {
	Route() martini.Router
}
