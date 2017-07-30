package elasticsearch

import (
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

// DownloadImage downloads image to filepath
func DownloadImage(url, filepath string) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Error(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{
		"status":   resp.Status,
		"size":     resp.ContentLength,
		"filepath": filepath,
	}).Debug("downloading file")

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error(err)
	}
}
