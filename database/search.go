package database

import (
	"math/rand"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/blevesearch/bleve"
)

const size = 100

func init() {
	rand.Seed(time.Now().Unix())
}

// SearchImage searches imagess by text/title and returns a random image
func (db *Database) SearchImage(search []string, itype string) (ImageMetaData, error) {

	var image ImageMetaData
	var images []ImageMetaData

	searchStr := strings.Join(removeNonAlphaNumericChars(search), " ")

	query := bleve.NewFuzzyQuery(searchStr)
	// query.SetFuzziness(1)
	// query := bleve.NewMatchPhraseQuery(searchStr)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 250
	// searchRequest.Fields
	// searchRequest.Highlight = bleve.NewHighlightWithStyle("ansi")
	searchResults, err := db.IDX.Search(searchRequest)
	if err != nil {
		return ImageMetaData{}, err
	}

	if searchResults.Total > 0 {
		for _, hit := range searchResults.Hits {
			doc, err := db.IDX.Document(hit.ID)
			if err != nil {
				return ImageMetaData{}, err
			}
			// log.Debugf("result doc string: %s", doc.GoString())
			for _, field := range doc.Fields {
				if field.Name() == "source" && string(field.Value()) == itype {
					image = ImageMetaData{}
					db.SQL.Find(&image, ImageMetaData{ID: hit.ID, Source: itype})
					images = append(images, image)
				}
			}
		}
		log.Debugf("search term (%s) returned %d results", searchStr, len(images))
		if len(images) > 0 {
			return images[rand.Intn(len(images))], nil
		}
	}
	if strings.EqualFold(itype, "xkcd") {
		return ImageMetaData{
			Title: "not found",
			Text:  searchStr,
			Path:  "images/default/xkcd.png"}, nil
	}
	if strings.EqualFold(itype, "giphy") {
		return ImageMetaData{Path: "images/default/giphy.gif"}, nil
	}
	return ImageMetaData{}, ErrNoImagesFound
}

// SearchGetAll searches imagess by text/title and returns all matching images
func (db *Database) SearchGetAll(search []string, itype string) ([]ImageMetaData, error) {

	var image ImageMetaData
	var images []ImageMetaData

	searchStr := strings.Join(removeNonAlphaNumericChars(search), " ")

	searchRequest := bleve.NewSearchRequest(bleve.NewMatchPhraseQuery(searchStr))
	searchRequest.Size = size

	searchResults, err := db.IDX.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if searchResults.Total > 0 {
		for _, hit := range searchResults.Hits {
			doc, err := db.IDX.Document(hit.ID)
			if err != nil {
				return nil, err
			}
			// log.Debugf("result doc string: %s", doc.GoString())
			for _, field := range doc.Fields {
				if field.Name() == "source" && string(field.Value()) == itype {
					image = ImageMetaData{}
					db.SQL.Find(&image, ImageMetaData{ID: hit.ID, Source: itype})
					images = append(images, image)
				}
			}
		}
		log.Debugf("search term (%s) returned %d results", searchStr, len(images))
		if len(images) > 0 {
			return images, nil
		}
	}

	if strings.EqualFold(itype, "xkcd") {
		return []ImageMetaData{ImageMetaData{
			Title: "not found",
			Text:  searchStr,
			Path:  "images/default/xkcd.png"}}, nil
	}
	if strings.EqualFold(itype, "giphy") {
		return []ImageMetaData{ImageMetaData{Path: "images/default/giphy.gif"}}, nil
	}
	return nil, ErrNoImagesFound
}

// SearchASCII searches ascii by keywords and returns a random matching ascii
func (db *Database) SearchASCII(keywords []string) (ASCIIData, error) {

	var ascii ASCIIData
	var asciis []ASCIIData

	keywordsStr := strings.Join(removeNonAlphaNumericChars(keywords), " ")

	query := bleve.NewMatchPhraseQuery(keywordsStr)
	searchRequest := bleve.NewSearchRequest(query)

	searchResults, err := db.IDX.Search(searchRequest)
	if err != nil {
		return ASCIIData{}, err
	}

	if searchResults.Total > 0 {
		for _, hit := range searchResults.Hits {
			doc, err := db.IDX.Document(hit.ID)
			if err != nil {
				return ASCIIData{}, err
			}
			// log.Debugf("result doc string: %s", doc.GoString())
			for _, field := range doc.Fields {
				if field.Name() == "source" && string(field.Value()) == "ascii" {
					ascii = ASCIIData{}
					db.SQL.Find(&ascii, ASCIIData{ID: hit.ID, Source: "ascii"})
					asciis = append(asciis, ascii)
				}
			}
		}
		log.Debugf("search term (%s) returned %d results", keywordsStr, len(asciis))
		if len(asciis) > 0 {
			return asciis[rand.Intn(len(asciis))], nil
		}
	}
	// return default 404 asciis
	return ASCIIData{
		ID:       "not found",
		Keywords: "10",
		Emoji:    "¯\\_(ツ)_/¯"}, nil
}
