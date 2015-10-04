package svcreg

import "github.com/enzian/go-msf/common"

// The GenericServiceRegistry manages service definitions and service versions in memory
type GenericServiceRegistry struct {
	services []common.ServiceDefinition
}

// CreateServiceDefinition create a new service definition entry
func (gsr GenericServiceRegistry) CreateServiceDefinition(identifierCode string, uriPrefix string, displayName string) (common.ServiceDefinition, error) {
	return common.ServiceDefinition{}, nil
}

// CreateServiceVersion attaches a new version to an existing service
func (gsr GenericServiceRegistry) CreateServiceVersion(identifierCode string, version string) (common.ServiceVersion, error) {
	return common.ServiceVersion{}, nil
}
