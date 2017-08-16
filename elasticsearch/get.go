package elasticsearch

import (
	"context"
	"fmt"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v5"
)

// GetImageByID returns the path to an image by id
func GetImageByID(id string) (ImageMetaData, error) {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return ImageMetaData{}, err
	}

	termQuery := elastic.NewTermQuery("id", id)
	searchResult, err := client.Search().
		Index("scifgif"). // search in index "twitter"
		Query(termQuery). // specify the query
		Size(1).          // take documents 0-9
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
	return ImageMetaData{}, fmt.Errorf("image id %s not found", id)
}
