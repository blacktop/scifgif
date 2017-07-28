package giphy

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/html"
)

const (
	// NumberOfGifs Total number of gifs to download
	NumberOfGifs = 1000
	giphyFolder  = "images/giphy"
)

func init() {
	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

// Helper function to pull the tag attribute from a Token
func getTag(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "tag"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

// Extract all tags from a given webpage
func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getTag(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}
		}
	}
}

func downloadImage(url string, filename string) {

	filepath := filepath.Join(giphyFolder, path.Base(filename+".gif"))

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

// GetAllGiphy havest all teh gifts
func GetAllGiphy() error {
	giphy := NewClient()

	for i := 0; i < NumberOfGifs; i += 25 {
		dataSearch, err := giphy.Search([]string{"reactions"}, i)
		// dataTrending, err := giphy.GetTrending()
		if err != nil {
			return err
		}
		for _, gif := range dataSearch.Data {
			downloadImage(gif.Images.Downsized.URL, gif.Slug)
			// fmt.Printf("GIPHY tags: %+v\n", gif.Tags)
		}
	}

	// foundUrls := make(map[string]bool)
	// seedUrls := os.Args[1:]

	// // Channels
	// chUrls := make(chan string)
	// chFinished := make(chan bool)

	// // Kick off the crawl process (concurrently)
	// for _, url := range seedUrls {
	// 	go crawl(url, chUrls, chFinished)
	// }

	// // Subscribe to both channels
	// for c := 0; c < len(seedUrls); {
	// 	select {
	// 	case url := <-chUrls:
	// 		foundUrls[url] = true
	// 	case <-chFinished:
	// 		c++
	// 	}
	// }

	// // We're done! Print the results...
	// fmt.Println("Found", len(foundUrls), "unique urls:")

	// for url := range foundUrls {
	// 	fmt.Println(" - " + url)
	// }

	// close(chUrls)
	return nil
}
