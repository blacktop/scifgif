package elasticsearch

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

// Size is the number of document results to return
const Size = 20

// SearchImage searches imagess by text/title and returns a random image
func SearchImage(search []string, itype string) (ImageMetaData, error) {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return ImageMetaData{}, err
	}

	searchStr := strings.Join(search, " ")

	// build randomly sorted search query
	q := elastic.NewMultiMatchQuery(searchStr, "title", "text").Operator("and") //.TieBreaker(0.3)
	// Search with a term query
	searchResult, err := client.Search().
		Index("scifgif"). // search in index "scifgif"
		Type(itype).      // only search supplied type images
		Query(q).         // specify the query
		Size(Size).
		Do(ctx) // execute
	if err != nil {
		return ImageMetaData{}, err
	}

	if searchResult.TotalHits() > 0 {
		var ityp ImageMetaData
		randomResult := rand.Intn(int(searchResult.TotalHits())) % Size
		for iter, item := range searchResult.Each(reflect.TypeOf(ityp)) {
			if i, ok := item.(ImageMetaData); ok {
				// return random image
				if iter == randomResult {
					log.WithFields(log.Fields{
						"total_hits":  searchResult.TotalHits(),
						"search_term": searchStr,
						"text":        i.Text,
					}).Debug("search found image")

					return i, nil
				}
			}
		}
	}

	log.WithFields(log.Fields{
		"type":        itype,
		"search_term": searchStr,
	}).Error("no found image")
	// return default 404 images
	if strings.EqualFold(itype, "xkcd") {
		return ImageMetaData{
			Title: "not found",
			Text:  searchStr,
			Path:  "images/default/xkcd.png"}, nil
	}
	if strings.EqualFold(itype, "giphy") {
		return ImageMetaData{Path: "images/default/giphy.gif"}, nil
	}
	return ImageMetaData{}, errors.New("no images found")
}
