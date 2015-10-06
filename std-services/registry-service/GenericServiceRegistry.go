package svcreg

import (
	"github.com/enzian/go-msf/common"
	"sync"
)

// The GenericServiceRegistry manages service definitions and service versions in memory
type GenericServiceRegistry struct {
	servicesMutex *sync.RWMutex
	services      []common.ServiceDefinition
}

// NewGenericServiceRegistry initializes a new instance of GenericServiceRegistry
func NewGenericServiceRegistry() GenericServiceRegistry {
	var svcReg = GenericServiceRegistry{}
	svcReg.servicesMutex = new(sync.RWMutex)
	return svcReg
}

// CreateServiceDefinition create a new service definition entry
func (gsr GenericServiceRegistry) CreateServiceDefinition(identifierCode string, uriPrefix string, displayName string) (*common.ServiceDefinition, error) {

	return nil, nil
}

// CreateServiceVersion attaches a new version to an existing service
func (gsr GenericServiceRegistry) CreateServiceVersion(identifierCode string, version string) (*common.ServiceVersion, error) {
	return nil, nil
}
