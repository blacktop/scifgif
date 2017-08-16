package elasticsearch

import (
	"context"
	"errors"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v5"
)

// GetRandomImage returns a random image path from source (xkcd/giphy)
func GetRandomImage(source string) (ImageMetaData, error) {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return ImageMetaData{}, err
	}

	// build random query
	q := elastic.NewFunctionScoreQuery().
		Query(elastic.NewTermQuery("source", source)).
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
		return ImageMetaData{}, err
	}

	if searchResult.TotalHits() > 0 {
		var ityp ImageMetaData
		for _, item := range searchResult.Each(reflect.TypeOf(ityp)) {
			if i, ok := item.(ImageMetaData); ok {
				return i, nil
			}
		}
	}
	return ImageMetaData{}, errors.New("no images found")
}
