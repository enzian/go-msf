package apisrv

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"

	"github.com/enzian/go-msf/common"
)

// APIInformationService keeps and maintains information about the APIs active routes by passing events through projections to the routing cache
type APIInformationService struct {
	projectionLines    map[string]([]Projection)
	routeCache         RouteCache
	EventChannel       chan common.Event
	eventHandlerCancel chan bool
}

// NewAPIInformationService initializes a new instance of ApiInformationService
func NewAPIInformationService() (*APIInformationService, error) {
	var ais = new(APIInformationService)
	ais.projectionLines = make(map[string]([]Projection))
	ais.EventChannel = make(chan common.Event)
	ais.eventHandlerCancel = make(chan bool)
	ais.routeCache = RouteCache{}
	return ais, nil
}

// ClassicAPIInfoService creates a new setup for the ApiInformationService that has all necessary
// message handlers and to create and update entries for services, versions, api-levels and host lists.
func ClassicAPIInfoService() (*APIInformationService, error) {

	var apiSvc, err = NewAPIInformationService()

	if err != nil {
		return nil, fmt.Errorf("Cannot start Api Service: %s", err)
	}

	var directoryHandlers = NewDirectory()
	apiSvc.Use("SERVICE_ADDED", directoryHandlers.AddServicePrefix)
	apiSvc.Use("API_VERSION_ADD", directoryHandlers.AddAPIVersion)
	apiSvc.Use("HOST_ADD", directoryHandlers.AddHost)

	return apiSvc, nil
}

// Use attaches a new projection to the given event
func (apiSvc *APIInformationService) Use(event string, proj Projection) {
	apiSvc.projectionLines[event] = append(apiSvc.projectionLines[event], proj)
}

//
// // Handle runns the given event though the given handler for the Action specified in the event
// func (apiSvc *APIInformationService) Handle(event Event) {
// 	apiSvc.eventChannel <- event
// }

// Start beginns processing events
func (apiSvc *APIInformationService) Start() {
	go apiSvc.processEvent()
}

// Stop stops the processing of events
func (apiSvc *APIInformationService) Stop() {
	apiSvc.eventHandlerCancel <- true
}

func (apiSvc *APIInformationService) processEvent() {
	for {
		select {
		case evnt := <-apiSvc.EventChannel:

			for _, projection := range apiSvc.projectionLines[evnt.Action] {
				newCache, err := projection(apiSvc.routeCache, evnt)
				if err != nil {
					break
				}
				apiSvc.routeCache = newCache
			}
		case <-apiSvc.eventHandlerCancel:
			break
		}
	}
}

// ReqForwarder is the handle function that proxys request to instances in the service cloud according to the route table.
func (apiSvc *APIInformationService) ReqForwarder(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var routeRegex = regexp.MustCompile(`(?P<domain>.*)\/api\/v(?P<version>(\d*\.)*\d*)\/(?P<prefix>[^\/]*)\/(?P<suffix>.*)`)
	var reqUrl = r.URL.String()
	if routeRegex.MatchString(r.URL.String()) {
		fmt.Println(reqUrl)
		match := routeRegex.FindStringSubmatch(reqUrl)
		result := make(map[string]string)
		for i, name := range routeRegex.SubexpNames() {
			result[name] = match[i]
		}

		if result["version"] != "" || result["prefix"] != "" {
			var key = fmt.Sprintf("%s,%s", result["version"], result["prefix"])
			var routes = apiSvc.routeCache[key]
			if len(routes) > 0 {
				nodeUrl, err := url.ParseRequestURI(fmt.Sprintf("http://%s/%s", routes[0], result["suffix"]))
				if err == nil {
					elapsed := time.Since(start)
					fmt.Printf("It took %s to proxy %s\n", elapsed, reqUrl)
					var client = httputil.NewSingleHostReverseProxy(nodeUrl)
					client.ServeHTTP(w, r)
					return
				}
			}
		}
	}

	http.Error(w, "malformed reqest uri", http.StatusBadRequest)
}
