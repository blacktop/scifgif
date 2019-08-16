package dilbert

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"github.com/blacktop/scifgif/database"
	"github.com/jpillora/backoff"
)

var attempt int

// MaxAttempts max number of download attempts
const MaxAttempts = 20

// Comic is the dilbert comic strip meta data
type Comic struct {
	Title      string
	Tags       []string
	ImageURL   string
	Transcript string
}

// GetComicMetaData gets all the comic strips meta data
func GetComicMetaData(url, date string, b *backoff.Backoff) Comic {
	comic := Comic{}

	if attempt > MaxAttempts {
		return comic
	}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.WithError(err).Error("goquery.NewDocument failed")
		// backoff and try again
		backoff := b.Duration()
		log.WithFields(log.Fields{
			"backoff": backoff,
			"attempt": attempt,
		}).Info("waiting to try to again")
		time.Sleep(backoff)
		// retry url meta data parse
		attempt++
		GetComicMetaData(url, date, b)
	}
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

	return comic
}

// GetAllDilbert havest all teh comics strips
func GetAllDilbert(folder string, date string) error {

	// open database
	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	count := 0
	attempt = 0

	b := &backoff.Backoff{
		//These are the defaults
		Min:    100 * time.Millisecond,
		Max:    1200 * time.Second,
		Factor: 3,
		Jitter: true,
	}

	if len(date) < 1 {
		// date = "1989-04-17"
		date = "2017-09-08"
	}
	start, _ := time.Parse("2006-01-02", date)

	for d := start; time.Now().After(d); d = d.AddDate(0, 0, 1) {
		date := fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
		dilbertURL := "http://dilbert.com/strip/" + date
		comic := GetComicMetaData(dilbertURL, date, b)

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
		b.Reset()
	}

	log.WithFields(log.Fields{"count": count}).Info("dilbert comic complete")
	return nil
}
