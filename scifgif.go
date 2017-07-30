package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
)

var (
	// Version stores the plugin's version
	Version string
	// BuildTime stores the plugin's build time
	BuildTime string
	// Token stores the webhook api token
	Token string
	// Host microservice host
	Host string
	// Port microservice port
	Port string
	// APIkey stores Giphy's API key
	APIkey string
)

// WebHookResponse mattermost webhook response struct
type WebHookResponse struct {
	Text     string `json:"text,omitempty"`
	Username string `json:"username,omitempty"`
	IconURL  string `json:"icon_url,omitempty"`
}

// getRandomXKCD serves a random xkcd comic
func getRandomXKCD(w http.ResponseWriter, r *http.Request) {
	path, err := elasticsearch.GetRandomImage("xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the xkcd source")
		log.Error(err)
		return
	}
	log.Debug(path)
	http.ServeFile(w, r, path)
}

// getXkcdByNumber serves a comic by it's number
func getXkcdByNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path, err := elasticsearch.GetImageByID(vars["number"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	http.ServeFile(w, r, path)
}

// getSearchXKCD serves a comic by searching for text
func getSearchXKCD(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path, err := elasticsearch.SearchImages("xkcd", r.Form["query"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	http.ServeFile(w, r, path)
}

// postXkcdMattermost handles xkcd webhook POST
func postXkcdMattermost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !strings.EqualFold(Token, r.Form["token"][0]) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, errors.New("unauthorized - bad token"))
		log.Error(errors.New("unauthorized - bad token"))
		return
	}

	path, err := elasticsearch.SearchImages("xkcd", r.Form["trigger_word"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		Text:     fmt.Sprintf("http://%s:%s/%s", Host, Port, path),
		Username: "scifgif",
		IconURL:  fmt.Sprintf("http://%s:%s/icon", Host, Port),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// getRandomGiphy serves a random giphy gif
func getRandomGiphy(w http.ResponseWriter, r *http.Request) {
	path, err := elasticsearch.GetRandomImage("giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the giphy source")
		log.Error(err)
		return
	}
	log.Debug(path)
	http.ServeFile(w, r, path)
}

// getSearchGiphy serves a comic by searching for text
func getSearchGiphy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path, err := elasticsearch.SearchImages("giphy", r.Form["query"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	http.ServeFile(w, r, path)
}

// postGiphyMattermost handles Giphy webhook POST
func postGiphyMattermost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !strings.EqualFold(Token, r.Form["token"][0]) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, errors.New("unauthorized - bad token"))
		log.Error(errors.New("unauthorized - bad token"))
		return
	}

	path, err := elasticsearch.SearchImages("giphy", r.Form["trigger_word"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		Text:     fmt.Sprintf("http://%s:%s/%s", Host, Port, path),
		Username: "scifgif",
		IconURL:  fmt.Sprintf("http://%s:%s/icon", Host, Port),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// getIcon serves scifgif icon
func getIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "images/icon.png")
}

// getImage serves scifgif icon
func getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["source"]
	file := vars["file"]

	// protect against directory traversal
	file = filepath.Clean(filepath.Base(file))
	log.Infof("GET images/%s/%s", folder, file)

	if _, err := os.Stat(filepath.Join("images", folder, file)); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "image not found")
		log.Error(err)
		return
	}

	http.ServeFile(w, r, filepath.Join("images", folder, file))
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
		cli.IntFlag{
			Name:   "number, N",
			Value:  1000,
			Usage:  "number of gifs to download",
			EnvVar: "IMAGE_NUMBER",
		},
		cli.StringFlag{
			Name:        "host",
			Value:       "",
			Usage:       "microservice host",
			EnvVar:      "SCIFGIF_HOST",
			Destination: &Host,
		},
		cli.StringFlag{
			Name:        "port",
			Value:       "3993",
			Usage:       "microservice port",
			EnvVar:      "SCIFGIF_PORT",
			Destination: &Port,
		},
		cli.StringFlag{
			Name:        "token",
			Value:       "",
			Usage:       "webhook token",
			EnvVar:      "SCIFGIF_TOKEN",
			Destination: &Token,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update images",
			Action: func(c *cli.Context) error {
				// start elasticsearch database
				elasticsearch.StartElasticsearch()
				// wait for elasticsearch to load
				err := elasticsearch.WaitForConnection(context.Background(), 60)
				if err != nil {
					log.Fatal(err)
				}

				// download Giphy gifs and ingest metadata into elasticsearch
				err = giphy.GetAllGiphy(giphyFolder, []string{"reactions"}, c.GlobalInt("number"))
				if err != nil {
					return err
				}
				// download xkcd comics and ingest metadata into elasticsearch
				err = xkcd.GetAllXkcd(xkcdFolder, c.GlobalInt("number"))
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

		if len(c.String("token")) == 0 {
			log.Warn("no webhook token set: --token")
		}

		// start elasticsearch database
		elasticsearch.StartElasticsearch()

		// wait for elasticsearch to load
		err := elasticsearch.WaitForConnection(context.Background(), c.Int("timeout"))
		if err != nil {
			log.Fatal(err)
		}

		// start web service
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/icon", getIcon).Methods("GET")
		router.HandleFunc("/images/{source}/{file}", getImage).Methods("GET")
		// xkcd routes
		router.HandleFunc("/xkcd", getRandomXKCD).Methods("GET")
		router.HandleFunc("/xkcd/{number}", getXkcdByNumber).Methods("GET")
		router.HandleFunc("/xkcd/search", getSearchXKCD).Methods("GET")
		router.HandleFunc("/xkcd/new_post", postXkcdMattermost).Methods("POST")
		// Giphy routes
		router.HandleFunc("/giphy", getRandomGiphy).Methods("GET")
		router.HandleFunc("/giphy/search", getSearchGiphy).Methods("GET")
		router.HandleFunc("/giphy/new_post", postGiphyMattermost).Methods("POST")
		// start microservice
		log.Info("web service listening on port :3993")
		log.Fatal(http.ListenAndServe(":"+Port, router))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
