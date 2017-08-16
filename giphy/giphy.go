package giphy

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/scifgif/elasticsearch"
)

// Helper function to pull the tag attribute from a Token
func getTags(url string, search string) []string {
	var gifTags []string

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Error(err)
	}

	// Find the script items
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		// extract tags from script
		re := regexp.MustCompile("\"tags\":.*\"]")
		if re.MatchString(s.Text()) {
			tags := re.FindString(s.Text())
			tags = strings.TrimPrefix(tags, "\"tags\": [")
			tags = strings.TrimSuffix(tags, "]")
			tagArray := strings.Split(tags, ",")
			for _, tag := range tagArray {
				tag = strings.Trim(strings.TrimSpace(tag), "\"\"")
				if !strings.Contains(search, tag) || len(tagArray) == 1 {
					gifTags = append(gifTags, tag)
				}
			}
		}
	})

	return gifTags
}

// GetAllGiphy havest all teh gifts
func GetAllGiphy(folder string, search []string, count int) error {
	giphy := NewClient()

	for i := 0; i < count; i += 25 {
		gifSearch, err := giphy.Search(search, i)
		// dataTrending, err := giphy.GetTrending()
		if err != nil {
			return err
		}
		for iter, gif := range gifSearch.Data {
			// download gif
			gifName := fmt.Sprintf("%s%d.gif", strings.Join(search, "-"), iter+i)
			filepath := filepath.Join(folder, path.Base(gifName))
			go elasticsearch.DownloadImage(gif.Images.Downsized.URL, filepath)
			srchStrs := strings.Join(search, " ")
			// index into elasticsearch
			elasticsearch.WriteImageToDatabase(elasticsearch.ImageMetaData{
				Name:   gif.Slug,
				ID:     gif.ID,
				Source: "giphy",
				Title:  gif.Source,
				Text:   strings.Join(getTags(gif.URL, srchStrs), " "),
				Path:   filepath,
			}, "giphy")
		}
	}
	return nil
}
