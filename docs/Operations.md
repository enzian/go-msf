# Operating go-msf

Operating `go-msf` is a breeze! Having the `go-msf` source , you can use it to run one of the three go-msf components:

* Service-Directory
* API-Server
* Echo-service (very much optional, it's a goodie for testing your setup!)

These three components can be hosted independently in a distributed or in a compact setup. The distributed setup it designed for large scale test/production operations while the compact mode is designed to help you host service-directory and API server integrated into the same service instance for simplicity reasons while developing your services.

### Integrated/compact setup

To host `go-msf` in the compact setup just launch it like this:
```
$ ./go-msf compact
```

If port :80 is already bound on your machine, you can use the `-u` parameter to specify another hostname and port to bind to:

```
$ ./go-msf compact -u ":8080"
```
or
```
$ ./go-msf compact -u "HOSTNAME:8080"
```
At this point your compact setup won't respond to any request coming into the API since there are no services registered. You can test your setup anyway because `go-msf` comes equipped with an echo service which you can launch like this:

```
$ ./go-msf echo -service-directory "localhost:8080/service-directory"
```
The `-service-directory` (or the shorthand `-srv`) parameter is there to tell the service to announce itself to the service-directory and where that service-directory instance can be found! In the compact setup the service-directory REST endpoint is located under `http://hostname:port/service-directory/...` while the API server is hosted under `http://hostname:port/api/...`. So, do not forget the tell your echo service to look for it in the right path (as shown in the sample above)!

The echo service will read the ContentType and the Content of the HTTP-Request and just return them in the response.

### Distributed setup

The distributed setup will be a little easier, since the services will be hosted on different machines. Since this functionality is not yet implemented - please stand by for changes that are yet to come! It will require you to deploy a distributed messaging system, which we'll specify later (though it will most probably be [NSQ](https://github.com/nsqio/nsq))
