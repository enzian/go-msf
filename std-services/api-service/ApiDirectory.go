package apisrv

import "fmt"

// Directory is used by the handler functions to keep track of the api versions and service prefixes that are needed to create the routing caches
type Directory struct {
	servicePrefixes []string
	apiversions     []string
}

// NewDirectory creates a new Directory with the handlers for adding services api versions and hosts
func NewDirectory() Directory {
	return Directory{}
}

// AddServicePrefix adds the service with it's prefix to the directory
func (dir *Directory) AddServicePrefix(rc RouteCache, evnt Event) (RouteCache, error) {
	for _, svc := range dir.servicePrefixes {
		if evnt.Data["prefix"] == svc {
			return nil, fmt.Errorf("service with the given prefix already registered")
		}
	}

	dir.servicePrefixes = append(dir.servicePrefixes, evnt.Data["prefix"])
	fmt.Println(fmt.Sprintf("Added new Service: %#v", dir.servicePrefixes))
	return rc, nil
}

// AddAPIVersion adds a version to the directory
func (dir *Directory) AddAPIVersion(rc RouteCache, evnt Event) (RouteCache, error) {
	for _, svc := range dir.apiversions {
		if evnt.Data["version"] == svc {
			return nil, fmt.Errorf("api with the given version already registered")
		}
	}

	dir.apiversions = append(dir.apiversions, evnt.Data["version"])
	fmt.Println("Added new Version")
	fmt.Println(fmt.Sprintf("Added new Version: %#v", dir.apiversions))
	return rc, nil
}

// AddHost adds a host for the given API version and service prefix
func (dir *Directory) AddHost(rc RouteCache, evnt Event) (RouteCache, error) {
	fmt.Println("Adding a new Host")
	fmt.Println(fmt.Sprintf("Versions: %#v", dir.apiversions))
	fmt.Println(fmt.Sprintf("Service: %#v", dir.servicePrefixes))

	var cacheKey = fmt.Sprintf("%s,%s", evnt.Data["version"], evnt.Data["prefix"])
	fmt.Println(fmt.Sprintf("Cache is now: %#v", rc))
	var routes = rc[cacheKey]
	var newRouteCache = append(routes, evnt.Data["uri"])
	rc[cacheKey] = newRouteCache
	fmt.Println("Added new Host")
	return rc, nil
}
