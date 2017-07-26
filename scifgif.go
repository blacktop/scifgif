package main

import (
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	xkcd "github.com/nishanths/go-xkcd"
)

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
