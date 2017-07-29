package giphy

import (
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/scifgif/elasticsearch"
)

func init() {
	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

// Helper function to pull the tag attribute from a Token
func getTags(url string) []string {
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
				gifTags = append(gifTags, tag)
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
		for _, gif := range gifSearch.Data {
			// download gif
			filepath := filepath.Join(folder, path.Base(gif.Slug+".gif"))
			go elasticsearch.DownloadImage(gif.Images.Downsized.URL, filepath)

			// index into elasticsearch
			elasticsearch.WriteImageToDatabase(elasticsearch.ImageMetaData{
				Name:  gif.Slug,
				ID:    gif.ID,
				Title: gif.Source,
				Text:  strings.Join(getTags(gif.URL), " "),
				Path:  filepath,
			})
		}
	}
	return nil
}
