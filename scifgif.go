package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/maliceio/malice/utils"
	xkcd "github.com/nishanths/go-xkcd"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	xkcdFolder  = "images/xkcd"
	giphyFolder = "images/giphy"
	mapping     = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"tweet":{
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
)

var (
	// ElasticAddr elasticsearch address to user for connections
	ElasticAddr string
)

func init() {
	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("ELASTICSEARCH", "elasticsearch"))
		log.Debug("Using elasticsearch address: ", ElasticAddr)
	}
}

// ImageMetaData image meta-data object
type ImageMetaData struct {
	ID    string `json:"id"`
	Name  string
	Title string
	Path  string
}

// SearchImages searches elasticsearch for images
func SearchImages(query string) error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	if err != nil {
		return err
	}

	// Search with a term query
	termQuery := elastic.NewTermQuery("title", query)
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

func downloadImage(url string) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join("images/xkcd", path.Base(url)), contents, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllXkcd() {
	client := xkcd.NewClient()
	latest, err := client.Latest()
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= latest.Number; i++ {
		comic, err := client.Get(i)
		if err != nil {
			log.Fatal(err)
		}
		downloadImage(comic.ImageURL)
	}
}

func getXKCD(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["file"]
	path := filepath.Join("images/xkcd", file)
	log.Println(path)
	http.ServeFile(w, r, path)
}

func main() {

	// static := os.Getenv("STATIC_DIR")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/xkcd/{file}", getXKCD).Methods("GET")
	log.Info("web service listening on port :3993")
	log.Fatal(http.ListenAndServe(":3993", router))
}
