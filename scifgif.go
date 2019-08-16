package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/blacktop/scifgif/ascii"
	"github.com/blacktop/scifgif/database"
	"github.com/blacktop/scifgif/dilbert"
	"github.com/blacktop/scifgif/giphy"
	"github.com/blacktop/scifgif/xkcd"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

const (
	xkcdFolder    = "images/xkcd"
	giphyFolder   = "images/giphy"
	dilbertFolder = "images/dilbert"
	contribFolder = "images/contrib"
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

	db *database.Database
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

func init() {
	log.SetHandler(logfmt.Default)
}

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
			Usage:  "database timeout (in seconds)",
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
			Name:   "date",
			Value:  "",
			Usage:  "dilbert comic start-from date",
			EnvVar: "IMAGE_DILBERT_DATE",
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

				if c.GlobalBool("verbose") {
					log.SetLevel(log.DebugLevel)
				}

				log.WithFields(log.Fields{
					"search_for": "reactions",
					"number":     c.GlobalInt("number"),
				}).Info("download reactions Giphy gifs and ingest metadata into database")
				err := giphy.GetAllGiphy(giphyFolder, []string{"reactions"}, c.GlobalInt("number"))
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"search_for": "star wars",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download star wars Giphy gifs and ingest metadata into database")
				err = giphy.GetAllGiphy(giphyFolder, []string{"star", "wars"}, min(c.GlobalInt("number"), 500))
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"search_for": "futurama",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download futurama Giphy gifs and ingest metadata into database")
				err = giphy.GetAllGiphy(giphyFolder, []string{"futurama"}, min(c.GlobalInt("number"), 250))
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"search_for": "rick and morty",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download rick and morty Giphy gifs and ingest metadata into database")
				err = giphy.GetAllGiphy(giphyFolder, []string{"rick", "and", "morty"}, min(c.GlobalInt("number"), 250))
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"search_for": "office space",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download office space Giphy gifs and ingest metadata into database")
				err = giphy.GetAllGiphy(giphyFolder, []string{"office", "space"}, min(c.GlobalInt("number"), 250))
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"search_for": "it crowd",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download it crowd Giphy gifs and ingest metadata into database")
				err = giphy.GetAllGiphy(giphyFolder, []string{"it", "crowd"}, min(c.GlobalInt("number"), 250))
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"search_for": "the office",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download the office Giphy gifs and ingest metadata into database")
				err = giphy.GetAllGiphy(giphyFolder, []string{"the", "office"}, min(c.GlobalInt("number"), 250))
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"search_for": "hacker",
					"number":     min(c.GlobalInt("number"), 250),
				}).Info("download futurama Giphy gifs and ingest metadata into database")
				err = giphy.GetAllGiphy(giphyFolder, []string{"hacker"}, min(c.GlobalInt("number"), 250))
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"number": c.GlobalInt("xkcd-count"),
				}).Info("download xkcd comics and ingest metadata into database")
				err = xkcd.GetAllXkcd(xkcdFolder, c.GlobalInt("xkcd-count"))
				if err != nil {
					return err
				}

				log.Info("load all ascii-emojis into database")
				err = ascii.GetAllASCIIEmoji()
				if err != nil {
					return err
				}

				log.WithFields(log.Fields{
					"date": c.GlobalString("date"),
				}).Info("download dilbert comics and ingest metadata into database")
				err = dilbert.GetAllDilbert(dilbertFolder, c.GlobalString("date"))
				if err != nil {
					return err
				}

				log.Info("* finalize database db")
				err = database.Finalize()
				if err != nil {
					return err
				}

				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		var err error

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		if len(c.String("token")) == 0 {
			log.Warn("no webhook token set: --token")
		}

		db, err = database.Open()
		if err != nil {
			return err
		}
		defer db.Close()

		// create http routes
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/icon/xkcd", getXkcdIcon).Methods("GET")
		router.HandleFunc("/icon/giphy", getGiphyIcon).Methods("GET")
		router.HandleFunc("/icon/dilbert", getDilbertIcon).Methods("GET")
		router.HandleFunc("/images", addImage).Methods("PUT")
		router.HandleFunc("/images/{source:(?:giphy|xkcd|dilbert|default|contrib)}/{file}", updateImageKeywords).Methods("PATCH")
		router.HandleFunc("/images/{source:(?:giphy|xkcd|dilbert|default|contrib)}/{file}", getImage).Methods("GET")
		router.HandleFunc("/images/{source:(?:giphy|xkcd|dilbert|default|contrib)}/{file}", deleteImage).Methods("DELETE")
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
		// Ascii-emoji routes
		router.HandleFunc("/ascii", getRandomASCII).Methods("GET")
		router.HandleFunc("/ascii/search", getSearchASCII).Methods("GET")
		// router.HandleFunc("/ascii/new_post", postASCIIMattermost).Methods("POST")
		router.HandleFunc("/ascii/slash", postASCIIMattermostSlash).Methods("POST")
		// Dilbert routes
		router.HandleFunc("/dilbert", getRandomDilbert).Methods("GET")
		router.HandleFunc("/dilbert/date/{date}", getDilbertByDate).Methods("GET")
		router.HandleFunc("/dilbert/search", getSearchDilbert).Methods("GET")
		router.HandleFunc("/dilbert/new_post", postDilbertMattermost).Methods("POST")
		router.HandleFunc("/dilbert/slash", postDilbertMattermostSlash).Methods("POST")
		// reactJS Web App routes
		router.HandleFunc("/web/search", getWebSearch).Methods("GET")
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/build")))

		// start microservice
		log.WithFields(log.Fields{
			"host":  Host,
			"port":  Port,
			"token": Token,
		}).Info("web service listening")
		err = http.ListenAndServe(":"+Port, httplog.New(router))
		log.WithError(err).Fatal("error listening")

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal("running app failed")
	}
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

func getMD5(f io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// this is overkill to protext against malicious URL from being created
func makeURL(scheme, host, port, path string) string {
	var u string
	var base *url.URL

	p, err := url.Parse(path)
	if err != nil {
		log.WithError(err).Error("url parse failed")
	}

	// don't show port if it is the HTTP/HTTPS port
	if strings.EqualFold(port, "80") || strings.EqualFold(port, "443") {
		u = fmt.Sprintf("%s://%s", scheme, host)
		base, err = url.Parse(u)
		if err != nil {
			log.WithError(err).Error("url parse failed")
		}
		return base.ResolveReference(p).String()
	}

	u = fmt.Sprintf("%s://%s:%s", scheme, host, port)
	base, err = url.Parse(u)
	if err != nil {
		log.WithError(err).Error("url parse failed")
	}

	return base.ResolveReference(p).String()
}
