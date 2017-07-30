package elasticsearch

import (
	"context"
	"errors"
	"reflect"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"
)

// SearchImages searches imagess by source and text and returns a random image
func SearchImages(source string, search []string) (string, error) {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return "", err
	}

	searchStr := strings.Join(search, " ")

	// build random query
	q := elastic.NewFunctionScoreQuery().
		Query(elastic.NewTermQuery("source", source)).
		Add(elastic.NewTermQuery("text", searchStr), elastic.NewWeightFactorFunction(1.5)).
		AddScoreFunc(elastic.NewWeightFactorFunction(3)).
		AddScoreFunc(elastic.NewRandomFunction()).
		Boost(5).
		ScoreMode("multiply")
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
				return i.Path, nil
			}
		}
	}
	// TODO: return default image when nothing found
	return "", errors.New("no images found")
}
