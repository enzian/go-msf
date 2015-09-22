package main

import (
	"fmt"
	"html"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "MSF API Server"
	app.Usage = "The API Server dispatches requests coming in to the according hosts in the service cloud"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "URI, u",
			Value: "localhost:2017",
			Usage: "hostname and port to bind to",
		},
	}
	app.Action = runRegistryService
	app.Run(os.Args)
}

func runRegistryService(c *cli.Context) {
	http.HandleFunc("/", handleApiRequest)
	http.ListenAndServe("localhost:2020", nil)
}

func handleApiRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
