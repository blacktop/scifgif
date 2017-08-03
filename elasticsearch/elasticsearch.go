package elasticsearch

import (
	"context"
	"errors"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

const mapping = `
{
"settings":{
  "number_of_shards": 1,
  "number_of_replicas": 0,
	"analysis": {
		"analyzer": {
			"my_english_analyzer": {
				"type": "standard",
				"stopwords": "_english_"
			}
		}
	}
},
"mappings":{
  "image":{
    "properties":{
      "id":{
        "type":"keyword"
      },
			"source":{
        "type":"keyword"
      },
      "name":{
        "type":"keyword"
      },
      "title":{
        "type":"text",
				"analyzer": "my_english_analyzer",
        "store": true,
        "fielddata": true
      },
      "text":{
        "type":"text",
				"analyzer": "my_english_analyzer",
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
	ID      string                `json:"id,omitempty"`
	Source  string                `json:"source,omitempty"`
	Name    string                `json:"name,omitempty"`
	Title   string                `json:"title,omitempty"`
	Text    string                `json:"text,omitempty"`
	Path    string                `json:"path,omitempty"`
	Suggest *elastic.SuggestField `json:"suggest_field,omitempty"`
}

// StartElasticsearch starts the elasticsearch database
func StartElasticsearch() {
	cmd := exec.Command("/elastic-entrypoint.sh", "elasticsearch", "-p", "/tmp/epid")
	cmd.Start()
}

// TestConnection tests the ElasticSearch connection
func TestConnection() (bool, error) {
	var err error

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return false, err
	}

	// Ping the Elasticsearch server to get e.g. the version number
	log.Debug("* attempting to PING elasticsearch")
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(context.Background())
	if err != nil {
		return false, err
	}

	log.WithFields(log.Fields{
		"code":    code,
		"cluster": info.ClusterName,
		"version": info.Version.Number,
	}).Debug("* elasticSearch connection successful.")

	if code == 200 {
		return true, err
	}
	return false, err
}

// WaitForConnection waits for connection to Elasticsearch to be ready
func WaitForConnection(ctx context.Context, timeout int, verbose bool) error {
	var ready bool
	var connErr error
	secondsWaited := 0

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

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
				log.Debugf("* elasticsearch came online after %d seconds", secondsWaited)
				return connErr
			}
			secondsWaited++
			time.Sleep(1 * time.Second)
		}
	}
}

// WriteImageToDatabase upserts image metadata into Database
func WriteImageToDatabase(image ImageMetaData, itype string) error {
	var err error
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return err
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("scifgif").Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		// Create a new index.
		createIndex, cerr := client.CreateIndex("scifgif").BodyString(mapping).Do(ctx)
		if cerr != nil {
			return cerr
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			log.Error(errors.New("index scifgif creation was not acknowledged"))
		} else {
			log.WithFields(log.Fields{"index": "scifgif"}).Info("index created")
		}
	}

	put, err := client.Index().
		Index("scifgif").
		Type(itype).
		OpType("index").
		BodyJson(image).
		Do(ctx)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id":    put.Id,
		"index": put.Index,
		"type":  put.Type,
	}).Debug("indexed image")

	return nil
}

// Finalize makes index read only optimized
func Finalize() error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return err
	}
	// Flush to make sure the documents got written.
	_, err = client.Flush().Do(ctx)
	if err != nil {
		return err
	}

	_, err = client.Forcemerge("scifgif").MaxNumSegments(1).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
