package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v5"
)

// SearchImages searches elasticsearch for images
func SearchImages(query string) error {
	ctx := context.Background()

	client, err := elastic.NewClient()
	if err != nil {
		return err
	}

	// Search with a term query
	termQuery := elastic.NewQueryStringQuery(query)
	searchResult, err := client.Search().
		Index("scifgif").    // search in index "twitter"
		Query(termQuery).    // specify the query
		Sort("title", true). // sort by "user" field, ascending
		From(0).Size(10).    // take documents 0-9
		Pretty(true).        // pretty print request and response JSON
		Do(ctx)              // execute
	if err != nil {
		return err
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ityp ImageMetaData
	for _, item := range searchResult.Each(reflect.TypeOf(ityp)) {
		if i, ok := item.(ImageMetaData); ok {
			fmt.Printf("Image  %s: %s\n", i.Name, i.Path)
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var i ImageMetaData
			err := json.Unmarshal(*hit.Source, &i)
			if err != nil {
				return err
			}

			// Work with image
			fmt.Printf("Image  %s: %s\n", i.Name, i.Path)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}
	return nil
}
