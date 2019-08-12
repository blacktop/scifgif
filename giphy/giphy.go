package giphy

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
	"github.com/blacktop/scifgif/database"
)

// Helper function to pull the tag attribute from a Token
func getTags(url string, search string) []string {
	var gifTags []string

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.WithError(err).Error("new document failed")
	}

	// Find the script items
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		// extract tags from script
		re := regexp.MustCompile("\"tags\":.*\"]")
		if re.MatchString(s.Text()) {
			tags := re.FindString(s.Text())
			tags = strings.TrimPrefix(tags, "\"tags\": [")
			tagParts := strings.SplitAfterN(tags, "]", 2)
			tags = strings.TrimSuffix(tagParts[0], "]")
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

	// open database
	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	gClient := NewClient()

	for i := 0; i < count; i += 25 {
		gifSearch, err := gClient.Search(search, i)
		// dataTrending, err := giphy.GetTrending()
		if err != nil {
			return err
		}

		for iter, gif := range gifSearch.Data {
			// check for a valid download URL
			_, err := url.ParseRequestURI(gif.Images.Downsized.URL)
			if err != nil {
				log.WithError(err).Errorf("url parsing failed for: %s", gif.Images.Downsized.URL)
				continue
			}
			// download gif
			gifName := fmt.Sprintf("%s%d.gif", strings.Join(search, "-"), iter+i)
			filepath := filepath.Join(folder, path.Base(gifName))
			go database.DownloadImage(gif.Images.Downsized.URL, filepath)
			srchStrs := strings.Join(search, " ")
			// index into bleve database
			db.WriteImageToDatabase(database.ImageMetaData{
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
