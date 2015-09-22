package svcreg

import (
	"fmt"
	"net/http"

	"github.com/ahmetalpbalkan/go-linq"
	"github.com/enzian/go-msf/common"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

// ServiceRegistryService is the service stub for the service registry
type ServiceRegistryService struct {
	apiVersions []*common.APIVersion
	services    []*common.ServiceDefinition
}

// Route return the route handler for this
func (s ServiceRegistryService) Route() martini.Router {
	r := martini.NewRouter()
	r.Group("/apiversions/", func(r martini.Router) {
		r.Get("", s.getAPIVersions)
		r.Post("", binding.Bind(common.APIVersion{}), s.createAPIVersion)
	})

	r.Group("/apiversions/:version/", func(r martini.Router) {
		r.Get("", s.getAPIVersion)
		r.Get("services/", s.getAPIServices)
		r.Post("services/", binding.Bind(common.ServiceVersion{}), s.createAPIService)
	})

	r.Group("/services/", func(r martini.Router) {
		r.Get("", s.getServices)
		r.Post("", binding.Bind(common.ServiceDefinition{}), s.createService)
		r.Delete(":service/", s.dropService)
	})

	r.Group("/services/:service/versions", func(r martini.Router) {
		r.Get("", s.getServiceVersions)
		r.Post("", binding.Bind(common.ServiceVersion{}), s.createServiceVersion)
	})

	r.Group("/services/:service/:version/hosts/", func(r martini.Router) {
		r.Get("", s.getServiceHosts)
		r.Post("", binding.Bind(common.ServiceHost{}), s.attachServiceHost)
	})

	return r
}

// NewServiceRegistryStub returns a new instance of service registry stub
func NewServiceRegistryStub() ServiceRegistryService {
	return ServiceRegistryService{
		apiVersions: []*common.APIVersion{},
		services:    []*common.ServiceDefinition{},
	}
}

func (s *ServiceRegistryService) getServiceVersions(r render.Render, parms martini.Params) {
	svc, found, err := s.getServiceByIdentifer(parms["service"])
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find service %s", parms["service"]),
			UserMessage:      "Failed to find the service specified",
		})
		return
	}
	r.JSON(http.StatusOK, svc.ServiceVersions)
}

func (s *ServiceRegistryService) createServiceVersion(r render.Render, sv common.ServiceVersion, parms martini.Params) {
	svc, found, err := s.getServiceByIdentifer(parms["service"])
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find service %s", parms["service"]),
			UserMessage:      "Failed to find the service specified",
		})
		return
	}
	var newServiceVersion = new(common.ServiceVersion)
	newServiceVersion.ServiceIdentifier = svc.Identifier
	newServiceVersion.Version = sv.Version
	newServiceVersion.ServiceHosts = []common.ServiceHost{}
	svc.ServiceVersions = append(svc.ServiceVersions, newServiceVersion)
}

func (s *ServiceRegistryService) getAPIVersion(r render.Render, parms martini.Params) {
	r.JSON(http.StatusOK, []common.APIVersion{})
}

func (s *ServiceRegistryService) getAPIVersions(r render.Render) {
	r.JSON(http.StatusOK, s.apiVersions)
}

func (s *ServiceRegistryService) getAPIServices(r render.Render, parms martini.Params) {
	var api, found, err = s.getAPIByVersion(parms["version"])
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find api version %s", parms["version"]),
			UserMessage:      "Failed to find matching api version",
		})
	} else {
		r.JSON(http.StatusOK, api.ServiceVersions)
	}
}

func (s *ServiceRegistryService) createAPIService(r render.Render, sv common.ServiceVersion, parms martini.Params) {
	var api, found, err = s.getAPIByVersion(parms["version"])
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find api version %s", parms["version"]),
			UserMessage:      "Failed to find matching api version",
		})
		return
	}
	svc, found, err := s.getServiceByIdentifer(sv.ServiceIdentifier)
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find service %s", sv.ServiceIdentifier),
			UserMessage:      "Failed to find the service specified",
		})
		return
	}

	svcv, found, err := linq.From(svc.ServiceVersions).FirstBy(func(x linq.T) (bool, error) {
		return x.(*common.ServiceVersion).ServiceIdentifier == svc.Identifier && x.(*common.ServiceVersion).Version == sv.Version, nil
	})
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find service %s with version %s", sv.ServiceIdentifier, sv.Version),
			UserMessage:      "Failed to find the service version specified",
		})
		return
	}
	api.ServiceVersions = append(api.ServiceVersions, svcv.(*common.ServiceVersion))
}

func (s *ServiceRegistryService) createAPIVersion(r render.Render, apiV common.APIVersion) {
	var _, found, err = s.getAPIByVersion(apiV.Version)

	if found == true || err != nil {
		r.JSON(http.StatusNotAcceptable, common.APIError{
			DeveloperMessage: fmt.Sprintf("Attempt to recreate an entry for api version %s", apiV.Version),
			UserMessage:      "Failed to create new api version",
		})
	} else {
		s.apiVersions = append(s.apiVersions, &apiV)
	}
}

func (s *ServiceRegistryService) getServices(r render.Render, parms martini.Params) {
	r.JSON(http.StatusOK, s.services)
}

func (s *ServiceRegistryService) createService(sd common.ServiceDefinition, r render.Render) {
	byIdentifier := func(x linq.T) (bool, error) {
		return x.(common.ServiceDefinition).Identifier == sd.Identifier, nil
	}
	var _, found, err = linq.From(s.services).Where(byIdentifier).First()
	if found == true || err != nil {
		r.JSON(http.StatusNotAcceptable, common.APIError{
			DeveloperMessage: fmt.Sprintf("Attempt to recreate an entry for service definition %s", sd.Identifier),
			UserMessage:      "Failed to create new service definition",
		})
	} else {
		sd.ServiceVersions = []*common.ServiceVersion{}
		s.services = append(s.services, &sd)
	}
}

func (s *ServiceRegistryService) dropService(r render.Render) {
	r.JSON(http.StatusNotAcceptable, common.APIError{
		DeveloperMessage: "Attempt to delete a service definition denied.",
		UserMessage:      "Action not allowed",
	})
}

func (s *ServiceRegistryService) getServiceHosts(r render.Render, parms martini.Params) {
	svc, found, err := s.getServiceByIdentifer(parms["service"])
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find service %s", parms["service"]),
			UserMessage:      "Failed to find the service specified",
		})
		return
	}

	byVersion := func(x linq.T) (bool, error) {
		return x.(*common.ServiceVersion).Version == parms["version"], nil
	}
	svcv, found, err := linq.From(svc.ServiceVersions).FirstBy(byVersion)
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find service version %s", parms["version"]),
			UserMessage:      "Failed to find the service version specified",
		})
		return
	}

	r.JSON(http.StatusOK, svcv.(*common.ServiceVersion).ServiceHosts)
}

func (s *ServiceRegistryService) attachServiceHost(r render.Render, host common.ServiceHost, parms martini.Params) {
	svc, found, err := s.getServiceByIdentifer(parms["service"])
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find service %s", parms["service"]),
			UserMessage:      "Failed to find the service specified",
		})
		return
	}

	byVersion := func(x linq.T) (bool, error) {
		return x.(*common.ServiceVersion).Version == parms["version"], nil
	}
	svcv, found, err := linq.From(svc.ServiceVersions).FirstBy(byVersion)
	if !found || err != nil {
		r.JSON(http.StatusNotFound, common.APIError{
			DeveloperMessage: fmt.Sprintf("Cannot find service version %s", parms["version"]),
			UserMessage:      "Failed to find the service version specified",
		})
		return
	}
	svcv.(*common.ServiceVersion).ServiceHosts = append(svcv.(*common.ServiceVersion).ServiceHosts, host)
}

func (s *ServiceRegistryService) getAPIByVersion(version string) (*common.APIVersion, bool, error) {
	byVersion := func(x linq.T) (bool, error) {
		return x.(*common.APIVersion).Version == version, nil
	}
	var x, found, err = linq.From(s.apiVersions).Where(byVersion).First()
	if !found {
		return nil, found, err
	}
	return x.(*common.APIVersion), found, err
}

func (s *ServiceRegistryService) getServiceByIdentifer(identifier string) (*common.ServiceDefinition, bool, error) {
	byIdentifier := func(x linq.T) (bool, error) {
		return x.(*common.ServiceDefinition).Identifier == identifier, nil
	}
	var x, found, err = linq.From(s.services).Where(byIdentifier).First()
	return x.(*common.ServiceDefinition), found, err
}
