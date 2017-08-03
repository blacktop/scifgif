package elasticsearch

import (
	"context"
	"math/rand"
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// SearchImage searches imagess by text/title and returns a random image
func SearchImage(search []string, itype string) (string, error) {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return "", err
	}

	searchStr := strings.Join(search, " ")

	// build randomly sorted search query
	q := elastic.NewMultiMatchQuery(searchStr, "text", "title").Operator("and").TieBreaker(0.3)
	// Search with a term query
	searchResult, err := client.Search().
		Index("scifgif"). // search in index "scifgif"
		Type(itype).      // only search supplied type images
		Query(q).         // specify the query
		Do(ctx)           // execute
	if err != nil {
		return "", err
	}

	if searchResult.TotalHits() > 0 {
		var ityp ImageMetaData
		randomResult := rand.Int63n(searchResult.TotalHits())
		for iter, item := range searchResult.Each(reflect.TypeOf(ityp)) {
			if i, ok := item.(ImageMetaData); ok {
				log.Info(iter)
				log.Info(randomResult)
				// return random image
				if iter == int(randomResult) {
					log.Info(i.Path)
					log.WithFields(log.Fields{
						"search_term": searchStr,
						"text":        i.Text,
					}).Debug("search found image")

					return i.Path, nil
				}
			}
		}
	}
	// TODO: return default image when nothing found
	return "images/default/nope.gif", nil
	// return "", errors.New("no images found")
}
