package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/enzian/go-msf/common"
	"net/http"
)

type ServiceDirectoryClient struct {
	sdURI string
}

func NewDirectoryClient(uri string) *ServiceDirectoryClient {
	subject := new(ServiceDirectoryClient)
	subject.sdURI = uri
	return subject
}

func (s *ServiceDirectoryClient) GetOrCreateServiceDefinition(identifier string, prefix string, displayName string) (*common.ServiceDefinition, error) {
	var subject = common.ServiceDefinition{
		Identifier:  identifier,
		URIPrefix:   prefix,
		DisplayName: displayName,
	}
	var content, err = json.Marshal(subject)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/services/", s.sdURI), "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("FATAL, status code: %s received", resp.Status)
	}

	return nil, nil
}

func (s *ServiceDirectoryClient) GetOrCreateServiceVersion(serviceIdentifier string, version string) (*common.ServiceVersion, error) {
	var subject = common.ServiceVersion{
		Version:           version,
		ServiceIdentifier: serviceIdentifier,
	}
	var content, err = json.Marshal(subject)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Calling POST on %s \n", fmt.Sprintf("http://%s/services/%s/versions/", s.sdURI, serviceIdentifier))
	resp, err := http.Post(fmt.Sprintf("http://%s/services/%s/versions/", s.sdURI, serviceIdentifier), "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FATAL, status code: %s received", resp.Status)
	}

	var sv = new(common.ServiceVersion)
	content = make([]byte, resp.ContentLength)
	json.Unmarshal(content, sv)
	return sv, nil
}

func (s *ServiceDirectoryClient) GetOrCreateApiVersion(version string, name string) (*common.APIVersion, error) {
	var subject = common.APIVersion{
		Version: version,
		Name:    name,
	}
	var content, err = json.Marshal(subject)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/apiversions/", s.sdURI), "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FATAL, status code: %s received", resp.Status)
	}

	var apiv = new(common.APIVersion)
	content = make([]byte, resp.ContentLength)
	json.Unmarshal(content, apiv)
	return apiv, nil
}

func (s *ServiceDirectoryClient) GetOrCreateServiceVersionApi(apiVersion string, serviceIdentifier string, serviceVersion string) error {
	var subject = common.ServiceVersion{
		Version:           serviceVersion,
		ServiceIdentifier: serviceIdentifier,
	}
	var content, err = json.Marshal(subject)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/apiversions/%s/services/", s.sdURI, apiVersion), "application/json", bytes.NewReader(content))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FATAL, status code: %s received", resp.Status)
	}

	return nil
}

func (s *ServiceDirectoryClient) AddServiceHost(serviceIdentifier string, serviceVersion string, hostUri string, hostState string) (*common.ServiceHost, error) {
	var subject = common.ServiceHost{
		URI:   hostUri,
		State: hostState,
	}
	var content, err = json.Marshal(subject)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/services/%s/%s/hosts/", s.sdURI, serviceIdentifier, serviceVersion), "application/json", bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FATAL, status code: %s received", resp.Status)
	}

	var host = new(common.ServiceHost)
	content = make([]byte, resp.ContentLength)
	json.Unmarshal(content, host)
	return host, nil
}
