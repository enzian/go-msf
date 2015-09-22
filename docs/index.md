# Microservices Framework - GO-MSF

The package go-msf ist aimed at making it easy to develop, test, run and monitor services in your service cloud independant of the technologies used in each service. It reduces the plumbing which usually has to be put in place to keep records of where and how to find and access a desired service.

Tough written in GO, GO-MSF does not require services to be written in GO. The interfaces to the service directories are specified in detail and can be used in any language/environment that can access a restfull service via HTTP.

It provides the two core components that are needed to run a distributed micro service architecture; The services registry and the API server.

You can find more information about the system under the following chapters:

* [Service Registry](ServiceRegistry.md "Service Registry")
* [API-Server](ApiServer.md "API-Server")
