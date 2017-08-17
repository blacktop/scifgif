package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/scifgif/elasticsearch"
	"github.com/blacktop/scifgif/giphy"
	"github.com/blacktop/scifgif/xkcd"
	"github.com/gorilla/handlers"
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
	Text         string `json:"text,omitempty"`
	Username     string `json:"username,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ResponseType string `json:"response_type,omitempty"`
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
			Value:  500,
			Usage:  "number of gifs to download",
			EnvVar: "IMAGE_NUMBER",
		},
		cli.IntFlag{
			Name:   "xkcd-count",
			Value:  -1,
			Usage:  "number of xkcd comics to download",
			EnvVar: "IMAGE_XKCD_COUNT",
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
				err := elasticsearch.WaitForConnection(context.Background(), 60, c.GlobalBool("verbose"))
				if err != nil {
					log.Fatal(err)
				}
				log.WithFields(log.Fields{
					"search_for": "reactions",
					"number":     c.GlobalInt("number"),
				}).Info("download Giphy gifs and ingest metadata into elasticsearch")
				err = giphy.GetAllGiphy(giphyFolder, []string{"reactions"}, c.GlobalInt("number"))
				if err != nil {
					return err
				}
				log.WithFields(log.Fields{
					"search_for": "star wars",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download star wars Giphy gifs and ingest metadata into elasticsearch")
				err = giphy.GetAllGiphy(giphyFolder, []string{"star", "wars"}, min(c.GlobalInt("number"), 500))
				if err != nil {
					return err
				}
				log.WithFields(log.Fields{
					"search_for": "futurama",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download futurama Giphy gifs and ingest metadata into elasticsearch")
				err = giphy.GetAllGiphy(giphyFolder, []string{"futurama"}, min(c.GlobalInt("number"), 500))
				if err != nil {
					return err
				}
				log.WithFields(log.Fields{
					"number": c.GlobalInt("number"),
				}).Info("download xkcd comics and ingest metadata into elasticsearch")
				err = xkcd.GetAllXkcd(xkcdFolder, c.GlobalInt("xkcd-count"))
				if err != nil {
					return err
				}
				log.Info("* finalize elasticsearch db")
				err = elasticsearch.Finalize()
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
		err := elasticsearch.WaitForConnection(context.Background(), c.Int("timeout"), c.Bool("verbose"))
		if err != nil {
			log.Fatal(err)
		}

		// create http routes
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/icon/xkcd", getXkcdIcon).Methods("GET")
		router.HandleFunc("/icon/giphy", getGiphyIcon).Methods("GET")
		router.HandleFunc("/images/{source:(?:giphy|xkcd|default)}/{file}", getImage).Methods("GET")
		router.HandleFunc("/images/{source:(?:giphy|xkcd|default)}/{file}", deleteImage).Methods("DELETE")
		// xkcd routes
		router.HandleFunc("/xkcd", getRandomXKCD).Methods("GET")
		router.HandleFunc("/xkcd/number/{number}", getXkcdByNumber).Methods("GET")
		router.HandleFunc("/xkcd/search", getSearchXKCD).Methods("GET")
		router.HandleFunc("/xkcd/new_post", postXkcdMattermost).Methods("POST")
		router.HandleFunc("/xkcd/slash", postXkcdMattermostSlash).Methods("POST")
		// Giphy routes
		router.HandleFunc("/giphy", getRandomGiphy).Methods("GET")
		router.HandleFunc("/giphy/search", getSearchGiphy).Methods("GET")
		router.HandleFunc("/giphy/new_post", postGiphyMattermost).Methods("POST")
		router.HandleFunc("/giphy/slash", postGiphyMattermostSlash).Methods("POST")
		// start microservice
		log.WithFields(log.Fields{
			"host":  Host,
			"port":  Port,
			"token": Token,
		}).Info("web service listening")
		loggedRouter := handlers.LoggingHandler(os.Stdout, router)
		log.Fatal(http.ListenAndServe(":"+Port, loggedRouter))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// getRandomXKCD serves a random xkcd comic
func getRandomXKCD(w http.ResponseWriter, r *http.Request) {
	image, err := elasticsearch.GetRandomImage("xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the xkcd source")
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getXkcdByNumber serves a comic by it's number
func getXkcdByNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	image, err := elasticsearch.GetImageByID(vars["number"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getSearchXKCD serves a comic by searching for text
func getSearchXKCD(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image, err := elasticsearch.SearchImage(r.Form["query"], "xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
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

	image, err := elasticsearch.SearchImage(r.Form["trigger_word"], "xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		Text: fmt.Sprintf("### %s\n>%s\n\n%s",
			image.Title,
			image.Text,
			makeURL("http", Host, Port, image.Path),
		),
		Username: "xkcd",
		IconURL:  makeURL("http", Host, Port, "icon/xkcd"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// postXkcdMattermostSlash handles xkcd webhook POST for use with a slash command
func postXkcdMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var image elasticsearch.ImageMetaData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error(errors.New("unauthorized - bad token"))
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random xkcd")
		image, err = elasticsearch.GetRandomImage("xkcd")
	} else if isNumeric(textArg) {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting xkcd by number")
		image, err = elasticsearch.GetImageByID(textArg)
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting xkcd by title")
		image, err = elasticsearch.SearchImage(r.Form["text"], "xkcd")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		ResponseType: "in_channel",
		Text: fmt.Sprintf("### %s\n>%s\n\non behalf of @%s %s",
			image.Title,
			image.Text,
			userName,
			makeURL("http", Host, Port, image.Path),
		),
		Username: "xkcd",
		IconURL:  makeURL("http", Host, Port, "icon/xkcd"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// getRandomGiphy serves a random giphy gif
func getRandomGiphy(w http.ResponseWriter, r *http.Request) {
	image, err := elasticsearch.GetRandomImage("giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the giphy source")
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getSearchGiphy serves a comic by searching for text
func getSearchGiphy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image, err := elasticsearch.SearchImage(r.Form["query"], "giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
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

	image, err := elasticsearch.SearchImage(r.Form["trigger_word"], "giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		Text:     makeURL("http", Host, Port, image.Path),
		Username: "scifgif",
		IconURL:  makeURL("http", Host, Port, "icon/giphy"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// postGiphyMattermostSlash handles giphy webhook POST for use with a slash command
func postGiphyMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var image elasticsearch.ImageMetaData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error(errors.New("unauthorized - bad token"))
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random gif")
		image, err = elasticsearch.GetRandomImage("giphy")
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting gif by keyword")
		image, err = elasticsearch.SearchImage(r.Form["text"], "giphy")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		ResponseType: "in_channel",
		Text: fmt.Sprintf("**%s** on behalf of @%s %s",
			textArg,
			userName,
			makeURL("http", Host, Port, image.Path),
		),
		Username: "scifgif",
		IconURL:  makeURL("http", Host, Port, "icon/giphy"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// getGiphyIcon serves giphy icon
func getGiphyIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "images/icons/giphy-icon.png")
}

// getXkcdIcon serves xkcd icon
func getXkcdIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "images/icons/xkcd-icon.jpg")
}

// getDefaultImage gets default image path
func getDefaultImage(source string) string {
	switch source {
	case "xkcd":
		return "images/default/xkcd.png"
	case "giphy":
		return "images/default/giphy.gif"
	default:
		return "images/default/giphy.gif"
	}
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["source"]
	file := vars["file"]

	// protect against directory traversal
	file = filepath.Clean(filepath.Base(file))
	path := filepath.Join("images", folder, file)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "image not found")
		log.Error(err)
		return
	}

	log.Infof("deleting images/%s/%s", folder, file)
	err := os.Remove(path)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		log.Error(err, "unable to remove image")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "image successfully removed")
		log.Error(err)
		return
	}
}

// getImage serves scifgif icon
func getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["source"]
	file := vars["file"]

	// protect against directory traversal
	file = filepath.Clean(filepath.Base(file))

	if _, err := os.Stat(filepath.Join("images", folder, file)); os.IsNotExist(err) {
		log.Debugf("GET default %s image", folder)
		http.ServeFile(w, r, getDefaultImage(folder))
		return
	}
	log.Debugf("GET images/%s/%s", folder, file)
	http.ServeFile(w, r, filepath.Join("images", folder, file))
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// this is overkill to protext against malicious URL from being created
func makeURL(scheme, host, port, path string) string {
	var u string
	var base *url.URL

	p, err := url.Parse(path)
	if err != nil {
		log.Error(err)
	}
	// don't show port if it is the HTTP/HTTPS port
	if strings.EqualFold(port, "80") || strings.EqualFold(port, "443") {
		u = fmt.Sprintf("%s://%s", scheme, host)
		base, err = url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}
		return base.ResolveReference(p).String()
	}

	u = fmt.Sprintf("%s://%s:%s", scheme, host, port)
	base, err = url.Parse(u)
	if err != nil {
		log.Error(err)
	}

	return base.ResolveReference(p).String()
}
