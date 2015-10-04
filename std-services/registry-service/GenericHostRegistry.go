package svcreg

import "github.com/enzian/go-msf/common"

// The GenericHostRegistry keeps an maintains a record of all service hosts currently active
type GenericHostRegistry struct {
	hosts []common.ServiceHost
}

func (ghr GenericHostRegistry) AddHost(serviceBaseURI string, serviceIdentifier string, serviceVersion string, state string) (common.ServiceHost, error) {
	return common.ServiceHost{}, nil
}
func (ghr GenericHostRegistry) SetHostState(hostID string, state string) error {
	return nil
}
