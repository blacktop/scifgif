package dilbert

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/scifgif/elasticsearch"
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
		log.Error(err)
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
	comic.Title = doc.Find(".comic-title-name").Text()
	// GET IMAGE URL
	doc.Find(".img-comic-container").Each(func(i int, s *goquery.Selection) {
		comic.ImageURL, _ = s.Find("img").Attr("src")
	})
	// GET TAGS
	doc.Find(".comic-tags").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(i int, a *goquery.Selection) {
			comic.Tags = append(comic.Tags, strings.TrimPrefix(a.Text(), "#"))
		})
	})
	// GET TRANSCRIPT
	id := "js-toggle-transcript-" + date
	transcripts := doc.Find("div#" + id).Text()
	// clean up string
	transcripts = strings.TrimSpace(transcripts)
	transcripts = strings.TrimPrefix(transcripts, "Transcript")
	transcripts = strings.TrimSpace(transcripts)
	comic.Transcript = transcripts

	return comic
}

// GetAllDilbert havest all teh comics strips
func GetAllDilbert(folder string, date string) error {
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
		url := "http://dilbert.com/strip/" + date
		comic := GetComicMetaData(url, date, b)

		if attempt > MaxAttempts {
			return errors.New("max number of attempts reached")
		}

		// download image
		log.WithFields(log.Fields{
			"id":    date,
			"title": comic.Title,
		}).Debug("downloading file")

		filepath := filepath.Join(folder, date+".jpg")
		go elasticsearch.DownloadImage(comic.ImageURL, filepath)

		// index into elasticsearch
		elasticsearch.WriteImageToDatabase(elasticsearch.ImageMetaData{
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
