package database

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/apex/log"
)

// DownloadImage downloads image to filepath
func DownloadImage(url, filepath string) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.WithError(err).Error("file create failed")
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
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
		log.WithError(err).Error("file copy failed")
	}
}

func removeNonAlphaNumericChars(searchTerms []string) []string {
	var cleaned []string
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, term := range searchTerms {
		processedString := reg.ReplaceAllString(term, " ")
		processedString = strings.TrimSpace(processedString)
		cleaned = append(cleaned, processedString)
	}
	return cleaned
}
