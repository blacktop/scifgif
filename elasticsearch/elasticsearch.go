package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"time"

	log "github.com/Sirupsen/logrus"
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
      "id":{
        "type":"keyword"
      },
      "name":{
        "type":"keyword"
      },
      "text":{
        "type":"text",
        "store": true,
        "fielddata": true
      },
      "text":{
        "type":"text",
        "store": true,
        "fielddata": true
      },
      "path":{
        "type":"keyword"
      },
      "suggest_field":{
        "type":"completion"
      }
    }
  }
}
}`

// ImageMetaData image meta-data object
type ImageMetaData struct {
	ID          string                `json:"id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Title       string                `json:"title,omitempty"`
	Text        string                `json:"text,omitempty"`
	Description string                `json:"description,omitempty"`
	Path        string                `json:"path,omitempty"`
	Suggest     *elastic.SuggestField `json:"suggest_field,omitempty"`
}

// StartElasticsearch starts the elasticsearch database
func StartElasticsearch() error {
	// _, err := utils.RunCommand(context.Background(), "/elastic-entrypoint.sh", "elasticsearch")
	// // log.Info(output)
	// return err
	cmd := exec.Command("/elastic-entrypoint.sh", "elasticsearch")
	cmd.Start()
	return nil
}

// TestConnection tests the ElasticSearch connection
func TestConnection() (bool, error) {

	var err error

	client, err := elastic.NewClient()
	if err != nil {
		return false, err
	}

	// Ping the Elasticsearch server to get e.g. the version number
	log.Debug("attempting to PING elasticsearch")
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(context.Background())
	if err != nil {
		return false, err
	}

	log.WithFields(log.Fields{
		"code":    code,
		"cluster": info.ClusterName,
		"version": info.Version.Number,
	}).Debug("elasticSearch connection successful.")

	if code == 200 {
		return true, err
	}
	return false, err
}

// WaitForConnection waits for connection to Elasticsearch to be ready
func WaitForConnection(ctx context.Context, timeout int) error {

	var ready bool
	var connErr error
	secondsWaited := 0

	connCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	log.Debug("===> trying to connect to elasticsearch")
	for {
		// Try to connect to Elasticsearch
		select {
		case <-connCtx.Done():
			log.WithFields(log.Fields{"timeout": timeout}).Error("connecting to elasticsearch timed out")
			return connErr
		default:
			ready, connErr = TestConnection()
			if ready {
				log.Infof("elasticsearch came online after %d seconds", secondsWaited)
				return connErr
			}
			secondsWaited++
			time.Sleep(1 * time.Second)
		}
	}
}

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

// WriteImageToDatabase upserts image metadata into Database
func WriteImageToDatabase(image ImageMetaData) error {
	var err error
	ctx := context.Background()

	client, err := elastic.NewClient()
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
	}).Debug("indexed image")

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("scifgif").Do(ctx)
	if err != nil {
		panic(err)
	}

	return err
}

// DownloadImage downloads image to filepath
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
