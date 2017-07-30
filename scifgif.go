package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/scifgif/elasticsearch"
	"github.com/blacktop/scifgif/giphy"
	"github.com/blacktop/scifgif/xkcd"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

const (
	xkcdFolder  = "images/xkcd"
	giphyFolder = "images/giphy"
	// NumberOfGifs Total number of gifs to download
	NumberOfGifs = 1000
)

var (
	// Version stores the plugin's version
	Version string
	// BuildTime stores the plugin's build time
	BuildTime string
	// ApiKey stores Giphy's API key
	// ApiKey string
)

// getRandomXKCD serves a random xkcd comic
func getRandomXKCD(w http.ResponseWriter, r *http.Request) {
	path, err := elasticsearch.GetRandomImage("xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the xkcd source")
		log.Error(err)
	}
	log.Debug(path)
	http.ServeFile(w, r, path)
}

// getXKCD is a request Hander
func getXKCD(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := filepath.Clean(filepath.Base(vars["file"]))
	path := filepath.Join(xkcdFolder, file)
	log.Println(path)
	http.ServeFile(w, r, path)
}

// getRandomGiphy serves a random giphy gif
func getRandomGiphy(w http.ResponseWriter, r *http.Request) {
	path, err := elasticsearch.GetRandomImage("giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the giphy source")
		log.Error(err)
	}
	log.Debug(path)
	http.ServeFile(w, r, path)
}

// getGiphy is a request Hander
func getGiphy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := filepath.Clean(filepath.Base(vars["file"]))
	path := filepath.Join(giphyFolder, file)
	log.Println(path)
	http.ServeFile(w, r, path)
}

var appHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

{{.Usage}}

Version: {{.Version}}{{if or .Author .Email}}

Author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

func main() {
	cli.AppHelpTemplate = appHelpTemplate
	app := cli.NewApp()

	app.Name = "scifgif"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Humorous Image Micro-Service"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.IntFlag{
			Name:   "timeout",
			Value:  60,
			Usage:  "elasticsearch timeout (in seconds)",
			EnvVar: "TIMEOUT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update images",
			Action: func(c *cli.Context) error {
				// start elasticsearch database
				err := elasticsearch.StartElasticsearch()
				if err != nil {
					return err
				}
				// wait for elasticsearch to load
				err = elasticsearch.WaitForConnection(context.Background(), 60)
				if err != nil {
					log.Fatal(err)
				}
				// download Giphy gifs and ingest metadata into elasticsearch
				err = giphy.GetAllGiphy(giphyFolder, []string{"reactions"}, NumberOfGifs)
				if err != nil {
					return err
				}
				// download xkcd comics and ingest metadata into elasticsearch
				err = xkcd.GetAllXkcd(xkcdFolder)
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		// start elasticsearch database
		err := elasticsearch.StartElasticsearch()
		if err != nil {
			return err
		}

		// wait for elasticsearch to load
		err = elasticsearch.WaitForConnection(context.Background(), c.Int("timeout"))
		if err != nil {
			log.Fatal(err)
		}

		// start web service
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/xkcd", getRandomXKCD).Methods("GET")
		router.HandleFunc("/xkcd/{file}", getXKCD).Methods("GET")
		router.HandleFunc("/giphy", getRandomGiphy).Methods("GET")
		router.HandleFunc("/giphy/{file}", getGiphy).Methods("GET")
		log.Info("web service listening on port :3993")
		log.Fatal(http.ListenAndServe(":3993", router))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
