package main

import (
	"fmt"
	"log"

	"github.com/nishanths/go-xkcd"
)

func main() {
	client := xkcd.NewClient()
	comic, err := client.Get(599)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s: %s", comic.Title, comic.ImageURL) // Apocalypse: http://imgs.xkcd.com/comics/apocalypse.png
}
