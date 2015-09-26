# Getting Started
To get started you will need the GO runtime and tools installed. If you haven't yet installed them please consult the GO projects [Getting started guide](https://golang.org/doc/install).

We also need you the install the dependency vendoring tool [`gvt`](https://github.com/FiloSottile/gvt). If you do not have `gvt` installed you can do this by simply typing the following on your command line of choice:
```
$ go get -u github.com/FiloSottile/gvt
```
If you have your `PATH` variable also pointing to the `/bin` directory in the `GOPATH` you will be able to use `gvt` from any directory. Check this by typing
```
$ gvt
gvt, a simple go vendoring tool based on gb-vendor.

Usage:
        gvt command [arguments]

The commands are:

        fetch       fetch a remote dependency
        rebuild     rebuild dependencies from manifest
        update      update a local dependency
        list        list dependencies one per line
        delete      delete a local dependency

Use "gvt help [command]" for more information about a command.
```

## Building and Running

First we'll get the the `go-msf` sources from Github:
```
$ go get https://github.com/enzian/go-msf
```
After that you should see the sources located unter `GOROOT/src/github.com/enzian/go-msf/`

Secondly we need to get you dependencies into the source tree. `gvt` will get the dependencies which are specified in the `/vendor/manifest` file and load the into them `/vendor/` folder for your. It also strips all git (or other VCS) metadata so you do not have to worry about that.

```
$ gvt rebuild
```

You are now ready to build the go-msf sources:
```
GOPATH/src/github.com/enizan/go-msf$ go build
```

This should produce a new executable in your source directory called `go-msf` or `go-msf` for the windows users in the world.

You can now start hosting the two core components of the microservice architecture. `go-msf` will help you with that! Just type `go-msf -help` and it will give you the options you have to host either the service-directory or the API servers. Use the provided parameters to link the services together when choosing to host them independently (which is something we seriously recommend for production environments!).

You can also host the service directory and the API server in a single service instance by executing `./go-msf.exe compact`.
