package database

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
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

func removeNonAlphaNumericChars(searchTerms []string) []string {
	var cleaned []string
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	for _, term := range searchTerms {
		processedString := reg.ReplaceAllString(term, " ")
		processedString = strings.TrimSpace(processedString)
		cleaned = append(cleaned, processedString)
	}
	return cleaned
}
