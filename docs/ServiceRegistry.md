# Service Registry

The service registry maintains the versions of the external API, the services that are active in the service cloud, the versions of these services and the service hosts that host a service at a specific version. It provides this information the the API servers and the service instances in the service cloud.

It exposes a restfull API that can be used by any instance in the service cloud as well as the API servers. It will persists configurations asynchronously to a document store and publish messages to a yet to be defined message bus to notify API-Server and other instances about changes to service, version and host lists in the directory.

The restfull API is documented swagger in the [Service Registry API yaml](service_registry_api_doc.yaml) file.

Sample:

POST http://localhost/service-directory/services/
{
"identifier":"md",
"uriprefix":"md",
"displayname":"Master Data Service"
}

POST http://localhost/service-directory/services/md/versions
{
"version":"1.0"
}

POST http://localhost/service-directory/apiversions/
{
"version":"1.0"
}

POST http://localhost/service-directory/apiversions/1.0/services/
{
"service":"md",
"version":"1.0"
}

POST http://localhost/service-directory/services/md/1.0/hosts/
{
"uri":"http://localhost:9090/someuri/",
"state":"healty"
}
