package elasticsearch

import (
	"context"
	"errors"
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// SearchImages searches imagess by source and text and returns a random image
func SearchImages(source string, search []string) (string, error) {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return "", err
	}

	searchStr := strings.Join(search, " ")

	// build randomly sorted search query
	q := elastic.NewFunctionScoreQuery().
		Query(elastic.NewTermQuery("text", searchStr)).
		AddScoreFunc(elastic.NewRandomFunction()).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	// Search with a term query
	searchResult, err := client.Search().
		Index("scifgif"). // search in index "scifgif"
		Query(q).         // specify the query
		Size(1).          // take single document
		Do(ctx)           // execute
	if err != nil {
		return "", err
	}

	if searchResult.TotalHits() > 0 {
		var ityp ImageMetaData
		for _, item := range searchResult.Each(reflect.TypeOf(ityp)) {
			if i, ok := item.(ImageMetaData); ok {

				log.WithFields(log.Fields{
					"search_term": searchStr,
					"text":        i.Text,
				}).Debug("search found image")

				return i.Path, nil
			}
		}
	}
	// TODO: return default image when nothing found
	return "", errors.New("no images found")
}
