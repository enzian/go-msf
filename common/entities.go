package common

// ServiceDefinition defines a services independant of it's version
type ServiceDefinition struct {
	HeadVersion     string            `json:"headversion,omitempty"`
	Identifier      string            `json:"identifier"`
	DisplayName     string            `json:"displayname"`
	URIPrefix       string            `json:"uriprefix"`
	ServiceVersions []*ServiceVersion `json:"versions"`
}

// APIVersion defines a revision of the entire API
type APIVersion struct {
	Version         string            `json:"version,omitempty"`
	Name            string            `json:"name,omitempty"`
	ServiceVersions []*ServiceVersion `json:"versions"`
}

// ServiceVersion defines the link between multiple versions of a service and the hosts that expose those services
type ServiceVersion struct {
	Version           string        `json:"version,omitempty"`
	ServiceIdentifier string        `json:"service"`
	ServiceHosts      []ServiceHost `json:"hosts"`
}

const (
	// HostStateActive represents the state of a host that can process requests as intended
	HostStateActive = "active"
	// HostStatePhaseOut represents a host that can process requests, but whant's to opt out of the service cloud
	HostStatePhaseOut = "phaseout"
	// HostStateInactive represents a host that will refuse to answer/process requests
	HostStateInactive = "inactive"
)

// ServiceHost defines a machine in which a Service can be found
type ServiceHost struct {
	URI   string `json:"uri"`
	State string `json:"state"`
}

// APIError is used to communicate errors between Services and to the user
type APIError struct {
	DeveloperMessage string `json:"developerMessage,omitempty"`
	UserMessage      string `json:"userMessage,omitempty"`
	ErrorCode        string `json:"errorCode,omitempty"`
	MoreInfo         string `json:"moreInfo,omitempty"`
}

// Event is a structre which represents a change to the route cache in some way.
type Event struct {
	Action string
	Data   map[string]string
}
