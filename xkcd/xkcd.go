package xkcd

import (
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/scifgif/elasticsearch"
	xkcd "github.com/nishanths/go-xkcd"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// GetAllXkcd havest all teh comics
func GetAllXkcd(folder string) error {
	client := xkcd.NewClient()
	latest, err := client.Latest()
	if err != nil {
		return err
	}

	log.Infof("there are %d xkcd comics availble", latest.Number)

	// get all images up to latest
	for i := 1; i <= latest.Number; i++ {
		comic, err := client.Get(i)
		if err != nil {
			log.Error(err)
			continue
		}
		basename := path.Base(comic.ImageURL)
		// download image
		log.WithFields(log.Fields{
			"id":    comic.Number,
			"title": comic.SafeTitle,
		}).Debug("downloading file")
		filepath := filepath.Join(folder, basename)
		go elasticsearch.DownloadImage(comic.ImageURL, filepath)

		var description string
		if len(comic.Transcript) == 0 {
			description = comic.Alt
		} else {
			description = comic.Transcript
		}
		// index into elasticsearch
		elasticsearch.WriteImageToDatabase(elasticsearch.ImageMetaData{
			Name:   strings.TrimSuffix(basename, path.Ext(basename)),
			ID:     string(comic.Number),
			Source: "xkcd",
			Title:  comic.Title,
			Text:   description,
			Path:   filepath,
		})
	}
	return nil
}
