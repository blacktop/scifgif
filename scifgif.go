package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	xkcd "github.com/blacktop/scifgif/xkcd"
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

// getXKCD is a request Hander
func getXKCD(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := filepath.Clean(filepath.Base(vars["file"]))
	path := filepath.Join(xkcdFolder, file)
	log.Println(path)
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
	}
	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update images",
			Action: func(c *cli.Context) error {
				// download Giphy gifs and ingest metadata into elasticsearch
				// err := giphy.GetAllGiphy(giphyFolder, NumberOfGifs)
				// if err != nil {
				// 	return err
				// }
				// download xkcd comics and ingest metadata into elasticsearch
				err := xkcd.GetAllXkcd(xkcdFolder)
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

		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/xkcd/{file}", getXKCD).Methods("GET")
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
