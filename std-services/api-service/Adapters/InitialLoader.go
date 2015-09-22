package apisrv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/enzian/go-msf/common"
)

// The InitialLoader reads all information from the registry service and publishes messages to populate the caches at startup
type InitialLoader struct {
	registryServiceURI string
	sendEvntChan       chan<- common.Event
}

// NewInitialLoader creates and initializes the new loader
func NewInitialLoader(URI string) InitialLoader {
	var loader = InitialLoader{}
	loader.registryServiceURI = URI
	return loader
}

// Start initializes the loader but does not read any data from the registry service yet.
func (loader *InitialLoader) Start(sendEventChan chan<- common.Event) {
	loader.sendEvntChan = sendEventChan
}

// Load ready all information from the registry service and published messages accordingly
func (loader *InitialLoader) Load() error {
	var serviceDefinitions, err = loader.loadServiceList()
	if err != nil {
		return err
	}

	for _, svd := range serviceDefinitions {
		var event = common.Event{
			Action: "SERVICE_ADDED",
			Data: map[string]string{
				"prefix": svd.URIPrefix,
			},
		}
		loader.sendEvntChan <- event
	}

	apiVersion, err := loader.loadApiVersionList()
	if err != nil {
		return err
	}
	for _, apiV := range apiVersion {
		var event = common.Event{
			Action: "API_VERSION_ADD",
			Data: map[string]string{
				"version": apiV.Version,
			},
		}
		loader.sendEvntChan <- event
	}
	return nil
}

func (loader *InitialLoader) loadServiceList() ([]common.ServiceDefinition, error) {
	var uri = fmt.Sprintf("http://%s/%s/", loader.registryServiceURI, "services")
	fmt.Println(fmt.Sprintf("Loading services from: %s", uri))

	var services []common.ServiceDefinition

	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &services)
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("Unmarshaled %v services", len(services)))
	return services, nil
}

func (loader *InitialLoader) loadApiVersionList() ([]common.APIVersion, error) {
	var uri = fmt.Sprintf("http://%s/%s/", loader.registryServiceURI, "apiversions")
	fmt.Println(fmt.Sprintf("Loading api versions from: %s", uri))

	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiVersions []common.APIVersion
	err = json.Unmarshal(data, &apiVersions)
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("Unmarshaled %v api versions", len(apiVersions)))
	return apiVersions, nil
}

// Stop would stop the adapter but since this adapter is invoked using Load(), this will not have any effect.
func (loader *InitialLoader) Stop() {

}
