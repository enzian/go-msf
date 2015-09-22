package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/enzian/go-msf/common"
	"github.com/enzian/go-msf/std-services/api-service"
	apiadapters "github.com/enzian/go-msf/std-services/api-service/Adapters"
	"github.com/enzian/go-msf/std-services/registry-service"
	"github.com/enzian/go-msf/vendor/github.com/martini-contrib/render"
	"github.com/go-martini/martini"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Name = "Microservice Tools"
	app.Usage = "hosts all tools necessary to run a distributed setup of service directories and API server for your msf environment."
	app.Commands = []cli.Command{
		{
			Name:    "compact",
			Aliases: []string{"compact"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "URI, u",
					Value: "localhost:80",
					Usage: "specify this parameter to override the endpoint the application will bind to.",
				},
			},
			Usage:  "creates a compact setup that co-hosts the API server alongside the service directory within the same process",
			Action: startCompact,
		},
		{
			Name:    "service-directory",
			Aliases: []string{""},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "URI, u",
					Value: "localhost:80",
					Usage: "specify this parameter to override the endpoint the application will bind to.",
				},
				cli.StringFlag{
					Name:  "nsq",
					Value: "localhost:4150",
					Usage: "specifies the nsq deamon to connect to.",
				},
			},
			Usage:  "creates an instance of the service directory",
			Action: startDirectory,
		},
		{
			Name:    "api",
			Aliases: []string{""},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "URI, u",
					Value: "localhost:80",
					Usage: "specify this parameter to override the endpoint the application will bind to.",
				},
				cli.StringFlag{
					Name:  "nsq",
					Value: "localhost:4150",
					Usage: "specifies the nsq deamon to connect to.",
				},
				cli.StringFlag{
					Name:  "service-directory, svr",
					Value: "",
					Usage: "specifies the service directorys URI",
				},
			},
			Usage:  "creates an instance of the api service.",
			Action: startAPI,
		},
	}
	app.Run(os.Args)
}

func startCompact(c *cli.Context) {
	fmt.Println("Starting compact service setup")

	var apiSvc, err = apisrv.ClassicAPIInfoService()
	if err != nil {
		fmt.Print(err)
		return
	}
	apiSvc.Start()

	http.HandleFunc("/api", myHandler)

	var m = martini.Classic()
	m.Use(render.Renderer())

	var dirSvc = svcreg.NewServiceRegistryStub(apiSvc.EventChannel)
	m.Action(dirSvc.Route().Handle)

	http.HandleFunc("/service-directory/", func(w http.ResponseWriter, r *http.Request) {
		r.URL, _ = url.Parse(strings.TrimPrefix(r.URL.String(), "/service-directory"))
		m.ServeHTTP(w, r)
	})

	go http.ListenAndServe(c.String("u"), nil)

	var initialLoader = apiadapters.NewInitialLoader(c.String("u") + "/service-directory")
	initialLoader.Start(apiSvc.EventChannel)

	err = initialLoader.Load()
	if err != nil {
		fmt.Print(err)
		return
	}

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	<-exitChan
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MyHandler Called")
}

func mySecHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MySecHandler Called")
}

func startAPI(c *cli.Context) {
	fmt.Println("Starting full-blown API server")
	var apiSvc, err = apisrv.ClassicAPIInfoService()

	if err != nil {
		fmt.Println(err)
		return
	}

	var initialLoader = apiadapters.NewInitialLoader(c.String("svr"))
	initialLoader.Start(apiSvc.EventChannel)

	apiSvc.Start()

	err = initialLoader.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}

func startDirectory(c *cli.Context) {
	fmt.Println("Starting full-blown service-directory")

	var m = martini.Classic()
	m.Use(render.Renderer())

	var dirSvc = svcreg.NewServiceRegistryStub(make(chan common.Event))
	m.Action(dirSvc.Route().Handle)

	http.HandleFunc("/service-directory/", func(w http.ResponseWriter, r *http.Request) {
		r.URL, _ = url.Parse(strings.TrimPrefix(r.URL.String(), "/service-directory"))
		m.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(c.String("u"), nil))
}
