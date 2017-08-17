package xkcd

import (
	"path"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/scifgif/elasticsearch"
	xkcd "github.com/nishanths/go-xkcd"
)

// GetAllXkcd havest all teh comics
func GetAllXkcd(folder string, count int) error {
	var start int

	client := xkcd.NewClient()
	latest, err := client.Latest()
	if err != nil {
		return err
	}

	log.Infof("there are %d xkcd comics availble", latest.Number)

	// only go back count number of comics from latest
	if (latest.Number-count) < 0 || count < 0 {
		start = 1
	} else {
		start = latest.Number - count
	}
	// get all images up to latest
	for i := start; i <= latest.Number; i++ {
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
			ID:     strconv.Itoa(comic.Number),
			Source: "xkcd",
			Title:  comic.Title,
			Text:   description,
			Path:   filepath,
		}, "xkcd")
	}
	return nil
}
