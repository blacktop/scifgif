package dilbert

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"github.com/blacktop/scifgif/database"
	"github.com/pkg/errors"
)

var (
	attempt int
	proxies []string
)

// MaxAttempts max number of download attempts
const MaxAttempts = 30

// Comic is the dilbert comic strip meta data
type Comic struct {
	Title      string
	Tags       []string
	ImageURL   string
	Transcript string
}

func init() {
	rand.Seed(time.Now().Unix())
}

func randomAgent() string {
	var userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	}
	return userAgents[rand.Int()%len(userAgents)]
}

func getRandomProxies() (string, error) {

	var proxy string

	if len(proxies) == 0 {
		doc, err := goquery.NewDocument("https://www.ip-adress.com/proxy-list")
		if err != nil {
			return "", err
		}

		doc.Find("table").Each(func(i int, tablehtml *goquery.Selection) {
			tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
				proxy = "http://" + rowhtml.Find("td").First().Text()
				if len(proxy) > 7 {
					proxies = append(proxies, proxy)
				}
			})
		})
	}

	return proxies[rand.Int()%len(proxies)], nil
}

// GetComicMetaData gets all the comic strips meta data
func GetComicMetaData(dilbertURL, date string) (Comic, error) {

	comic := Comic{}

	if attempt > MaxAttempts {
		return comic, fmt.Errorf("attempts exceeded max attempts of %d", MaxAttempts)
	}

	proxyURL, err := url.Parse(proxies[attempt])
	if err != nil {
		return Comic{}, errors.Wrap(err, "parsing proxy URL failed")
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   60 * time.Second,
			ResponseHeaderTimeout: 60 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 120 * time.Second,
	}

	req, _ := http.NewRequest("GET", dilbertURL, nil)
	req.Header.Set("User-Agent", randomAgent())

	res, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error("client Do request failed")
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.WithError(err).Error("goquery.NewDocument failed")
	}

	if doc != nil {
		// GET TITLE
		doc.Find(".comic-title-name").Each(func(i int, s *goquery.Selection) {
			comic.Title = s.Text()
		})

		// GET IMAGE URL
		doc.Find(".img-comic-container").Each(func(i int, s *goquery.Selection) {
			comic.ImageURL, _ = s.Find("img").Attr("src")
			comic.ImageURL = "http:" + comic.ImageURL
		})

		// GET TAGS
		doc.Find(".comic-tags").Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(i int, a *goquery.Selection) {
				comic.Tags = append(comic.Tags, strings.TrimPrefix(a.Text(), "#"))
			})
		})

		// GET TRANSCRIPT
		id := "js-toggle-transcript-" + date
		doc.Find("div#" + id).Each(func(i int, s *goquery.Selection) {
			comic.Transcript = strings.TrimSpace(s.Text())
			comic.Transcript = strings.TrimPrefix(comic.Transcript, "Transcript")
			comic.Transcript = strings.TrimSpace(comic.Transcript)
		})

		return comic, nil
	}

	attempt++
	log.WithFields(log.Fields{
		"attempt": attempt,
		"proxy":   proxies[attempt],
	}).Info("retrying again")
	// retry url meta data parse
	return GetComicMetaData(dilbertURL, date)
}

// GetAllDilbert havest all teh comics strips
func GetAllDilbert(folder string, date string) error {

	// open database
	db, err := database.Open()
	if err != nil {
		return errors.Wrap(err, "opening database failed")
	}
	defer db.Close()

	count := 0
	attempt = 0

	_, err = getRandomProxies()
	if err != nil {
		return errors.Wrap(err, "getting a random proxy URLs failed")
	}

	if len(date) < 1 {
		// date = "1989-04-17"
		date = "2018-01-01"
	}
	start, _ := time.Parse("2006-01-02", date)

	for d := start; time.Now().After(d); d = d.AddDate(0, 0, 1) {
		date := fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
		comic, err := GetComicMetaData("http://dilbert.com/strip/"+date, date)
		if err != nil {
			return errors.Wrap(err, "getting comic metadata failed")
		}

		// check for a valid download URL
		dlURL, err := url.ParseRequestURI(comic.ImageURL)
		if err != nil {
			log.WithError(err).Errorf("url parsing failed for: %s", comic.ImageURL)
			continue
		}

		if attempt > MaxAttempts {
			return errors.New("max number of attempts reached")
		}

		// download image
		log.WithFields(log.Fields{
			"id":    date,
			"title": comic.Title,
			"url":   dlURL.String(),
		}).Debug("downloading file")

		filepath := filepath.Join(folder, date+".jpg")
		go database.DownloadImage(dlURL.String(), filepath)

		// index into bleve database
		db.WriteImageToDatabase(database.ImageMetaData{
			Name:   comic.Title,
			ID:     date,
			Source: "dilbert",
			Title:  strings.Join(comic.Tags, " "),
			Text:   comic.Transcript,
			Path:   filepath,
		}, "dilbert")

		// incr count, reset attempts and reset backoff
		count++
		attempt = 0
	}

	log.WithFields(log.Fields{"count": count}).Info("dilbert comic complete")
	return nil
}
