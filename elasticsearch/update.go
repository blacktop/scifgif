package elasticsearch

import (
	"context"
	"errors"

	log "github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

// UpdateKeywords adds new keywords to an image's search text
func UpdateKeywords(image ImageMetaData) error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return err
	}

	termQuery := elastic.NewTermQuery("id", image.ID)
	searchResult, err := client.Search().
		Index("scifgif"). // search in index "twitter"
		Query(termQuery). // specify the query
		Size(1).          // take documents 0-9
		Do(ctx)           // execute
	if err != nil {
		return err
	}

	if searchResult.TotalHits() < 1 {
		return errors.New("update image not found")
	}

	hit := searchResult.Hits.Hits[0]
	updateResult, err := client.Update().
		Index(hit.Index).
		Type(hit.Type).
		Id(hit.Id).
		Doc(map[string]interface{}{"text": image.Text}).
		DetectNoop(true).
		Do(ctx) // execute
	if err != nil {
		return err
	}

	log.Debugln("%#v", updateResult)
	return nil
}
