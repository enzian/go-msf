package svcreg

import (
	"fmt"
	linq "github.com/ahmetalpbalkan/go-linq"
	"github.com/enzian/go-msf/common"
	"sync"
)

// The GenericServiceRegistry manages service definitions and service versions in memory
type GenericServiceRegistry struct {
	servicesMutex *sync.RWMutex
	services      []*common.ServiceDefinition
}

// NewGenericServiceRegistry initializes a new instance of GenericServiceRegistry
func NewGenericServiceRegistry() *GenericServiceRegistry {
	var svcReg = GenericServiceRegistry{}
	svcReg.servicesMutex = new(sync.RWMutex)
	svcReg.services = []*common.ServiceDefinition{}
	return &svcReg
}

// CreateServiceDefinition create a new service definition entry
func (gsr *GenericServiceRegistry) CreateServiceDefinition(identifierCode string, uriPrefix string, displayName string) (*common.ServiceDefinition, error) {
	var filterByID = func(t linq.T) (bool, error) {
		return t.(*common.ServiceDefinition).Identifier == identifierCode, nil
	}
	var filterByPrefix = func(t linq.T) (bool, error) { return t.(common.ServiceDefinition).URIPrefix == uriPrefix, nil }

	// Lock the service definitions slices write lock
	gsr.servicesMutex.Lock()

	// Try to find an existing service definition with the given identifier. If so - fail and return error.
	var found, err = linq.From(gsr.services).AnyWith(filterByID)
	if found {
		return nil, fmt.Errorf("Attempt at overriding service identifier: %s", identifierCode)
	} else if err != nil {
		return nil, fmt.Errorf("Suffered error while filtering for service with identifier: %s", identifierCode)
	}

	// Try to find an existing service definition with the given URI prefix. If so - fail and return error.
	found, err = linq.From(gsr.services).AnyWith(filterByPrefix)
	if found {
		return nil, fmt.Errorf("Attempt at overriding service URI prefix: %s", uriPrefix)
	} else if err != nil {
		return nil, fmt.Errorf("Suffered error while filtering for service with URI prefix: %s", uriPrefix)
	}

	var sd = common.ServiceDefinition{
		Identifier:  identifierCode,
		URIPrefix:   uriPrefix,
		DisplayName: displayName,
	}

	gsr.services = append(gsr.services, &sd)

	// Unlock the service definitions slices write lock
	gsr.servicesMutex.Unlock()

	return &sd, nil
}

// CreateServiceVersion attaches a new version to an existing service
func (gsr *GenericServiceRegistry) CreateServiceVersion(serviceIdentifier string, version string) (*common.ServiceVersion, error) {
	var filterByID = func(t linq.T) (bool, error) {
		return t.(*common.ServiceDefinition).Identifier == serviceIdentifier, nil
	}
	var versionByID = func(t linq.T) (bool, error) {
		return t.(common.ServiceVersion).Version == version, nil
	}

	var element, found, err = linq.From(gsr.services).FirstBy(filterByID)
	if !found {
		return nil, fmt.Errorf("Attempt at overriding service identifier: %s", serviceIdentifier)
	} else if err != nil {
		return nil, fmt.Errorf("Suffered error while filtering for service with identifier: %s", serviceIdentifier)
	}

	found, err = linq.From((element.(*common.ServiceDefinition)).ServiceVersions).AnyWith(versionByID)
	if found {
		return nil, fmt.Errorf("Attempt at overriding service %s version: %s", serviceIdentifier, version)
	} else if err != nil {
		return nil, fmt.Errorf("Suffered error while filtering for service version: %s", version)
	}

	var sv = common.ServiceVersion{
		Version:           version,
		ServiceIdentifier: serviceIdentifier,
	}

	var definition = (element.(*common.ServiceDefinition))
	definition.ServiceVersions = append(definition.ServiceVersions, sv)

	return &sv, nil
}
