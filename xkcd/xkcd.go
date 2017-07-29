package xkcd

import (
	"path"
	"path/filepath"

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

	log.Infof("there are %d xkcd comics availble\n", latest.Number)

	// get all images up to latest
	for i := 1; i <= latest.Number; i++ {
		comic, err := client.Get(i)
		if err != nil {
			log.Error(err)
			continue
		}

		// download image
		log.WithFields(log.Fields{
			"id":    comic.Number,
			"title": comic.SafeTitle,
		}).Debug("downloading file")
		filepath := filepath.Join(folder, path.Base(comic.ImageURL))
		go elasticsearch.DownloadImage(comic.ImageURL, filepath)

		// index into elasticsearch
		elasticsearch.WriteImageToDatabase(elasticsearch.ImageMetaData{
			Name: comic.Title,
			Path: filepath,
		})
	}
	return nil
}
