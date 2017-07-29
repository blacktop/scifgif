package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/utils"
	elastic "gopkg.in/olivere/elastic.v5"
)

const mapping = `
{
"settings":{
  "number_of_shards": 1,
  "number_of_replicas": 0
},
"mappings":{
  "image":{
    "properties":{
      "user":{
        "type":"keyword"
      },
      "message":{
        "type":"text",
        "store": true,
        "fielddata": true
      },
      "image":{
        "type":"keyword"
      },
      "created":{
        "type":"date"
      },
      "tags":{
        "type":"keyword"
      },
      "location":{
        "type":"geo_point"
      },
      "suggest_field":{
        "type":"completion"
      }
    }
  }
}
}`

// ElasticAddr elasticsearch address to user for connections
var ElasticAddr string

func init() {
	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("ELASTICSEARCH", "127.0.0.1"))
		log.Debug("Using elasticsearch address: ", ElasticAddr)
	}
}

// ImageMetaData image meta-data object
type ImageMetaData struct {
	ID      string                `json:"id,omitempty"`
	Name    string                `json:"name,omitempty"`
	Title   string                `json:"title,omitempty"`
	Path    string                `json:"path,omitempty"`
	Suggest *elastic.SuggestField `json:"suggest_field,omitempty"`
}

// SearchImages searches elasticsearch for images
func SearchImages(query string) error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
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

// WriteImageToDatabase upserts image metadata into Database
func WriteImageToDatabase(image ImageMetaData) error {
	var err error
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	if err != nil {
		return err
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("scifgif").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("scifgif").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	put, err := client.Index().
		Index("scifgif").
		Type("image").
		OpType("index").
		BodyJson(image).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	log.WithFields(log.Fields{
		"id":    put.Id,
		"index": put.Index,
		"type":  put.Type,
	}).Debug("Indexed image.")

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("scifgif").Do(ctx)
	if err != nil {
		panic(err)
	}

	return err
}

func DownloadImage(url, filepath string) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Error(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{
		"status":   resp.Status,
		"size":     resp.ContentLength,
		"filepath": filepath,
	}).Debug("downloading file")

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error(err)
	}
}
