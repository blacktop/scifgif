package database

import (
	"math/rand"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
)

const size = 75

func init() {
	rand.Seed(time.Now().Unix())
}

// SearchImage searches imagess by text/title and returns a random image
func (db *Database) SearchImage(search []string, itype string) (ImageMetaData, error) {

	searchStr := strings.Join(removeNonAlphaNumericChars(search), " ")

	// query := bleve.NewFuzzyQuery(searchStr)
	// query.SetFuzziness(1)
	query := bleve.NewMatchPhraseQuery(searchStr)
	searchRequest := bleve.NewSearchRequest(query)
	// searchRequest.Fields
	// searchRequest.Highlight = bleve.NewHighlightWithStyle("ansi")
	searchResults, err := db.IDX.Search(searchRequest)
	if err != nil {
		return ImageMetaData{}, err
	}

	if searchResults.Total > 0 {
		var image ImageMetaData
		// id := searchResults.Hits[0].ID
		id := searchResults.Hits[rand.Intn(len(searchResults.Hits))].ID

		if db.SQL.Find(&image, ImageMetaData{ID: id, Source: itype}).RecordNotFound() {
			return ImageMetaData{}, ErrNoImagesFound
		}

		return image, nil
	}

	return ImageMetaData{}, ErrNoImagesFound
}

// SearchGetAll searches imagess by text/title and returns all matching images
func (db *Database) SearchGetAll(search []string, itype string) ([]ImageMetaData, error) {

	var image ImageMetaData
	var images []ImageMetaData

	searchStr := strings.Join(removeNonAlphaNumericChars(search), " ")

	query := bleve.NewMatchPhraseQuery(searchStr)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = size
	searchResults, err := db.IDX.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if searchResults.Total > 0 {
		for _, hit := range searchResults.Hits {
			image = ImageMetaData{}
			db.SQL.Find(&image, ImageMetaData{ID: hit.ID, Source: itype})
			images = append(images, image)
		}

		return images, nil
	}

	return []ImageMetaData{}, ErrNoImagesFound
}

// SearchASCII searches ascii by keywords and returns a random matching ascii
func SearchASCII(keywords []string) (ASCIIData, error) {

	// keywordsStr := strings.Join(removeNonAlphaNumericChars(keywords), " ")

	// return ASCIIData{}, ErrNoASCIIFound

	// return default 404 images
	return ASCIIData{
		ID:       "not found",
		Keywords: "10",
		Emoji:    "¯\\_(ツ)_/¯"}, nil
}
