package common

import "github.com/go-martini/martini"

// A ServiceStub returns routing parameters for multiple services to live on the same endpoint
type ServiceStub interface {
	Route() martini.Router
}

// IServiceRegistry abstracts the management of service definitions and versions
type IServiceRegistry interface {
	CreateServiceDefinition(identifierCode string, uriPrefix string, displayName string) (*ServiceDefinition, error)
	CreateServiceVersion(identifierCode string, version string) (*ServiceVersion, error)
}

// IApiRegistry abstracts the management of api revisions
type IApiRegistry interface {
	CreateApiVersion(version string, displayName string) (*ServiceDefinition, error)
	LinkServiceVersion(apiVersion string, serviceIdentifier string, serviceVersion string) error
}

// IHostRegistry abstracts the management of hosts for given services and versions
type IHostRegistry interface {
	AddHost(serviceBaseURI string, serviceIdentifier string, serviceVersion string, state string) (*ServiceHost, error)
	SetHostState(hostID string, state string) error
}
