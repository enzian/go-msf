package svcreg

import "github.com/enzian/go-msf/common"

// The GenericAPIVersionRegistry keeps an maintains records of all active API versions.
type GenericAPIVersionRegistry struct {
	versions []common.APIVersion
}

// CreateAPIVersion create a new API version in the directory
func (gap GenericAPIVersionRegistry) CreateAPIVersion(version string, displayName string) (common.ServiceDefinition, error) {
	return common.ServiceDefinition{}, nil
}

// LinkServiceVersion links an existing service version to a given (also existing) api version
func (gap GenericAPIVersionRegistry) LinkServiceVersion(apiVersion string, serviceIdentifier string, serviceVersion string) error {
	return nil
}
