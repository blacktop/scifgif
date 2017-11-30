package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/blacktop/scifgif/elasticsearch"
	"github.com/gorilla/mux"
)

// const staticDir = web

// func staticHandler(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Cache-Control", "max-age=31556926, public")
// 		if strings.HasSuffix(r.URL.Path, "/") {
// 			http.NotFound(w, r)
// 			return
// 		}
// 		handler.ServeHTTP(w, r)
// 	})
// }

// func root(c *api.Context, w http.ResponseWriter, r *http.Request) {

// 	w.Header().Set("Cache-Control", "no-cache, max-age=31556926, public")

// 	staticDir, _ := utils.FindDir(model.CLIENT_DIR)
// 	http.ServeFile(w, r, staticDir+"root.html")
// }

// getRandomXKCD serves a random xkcd comic
func getRandomXKCD(w http.ResponseWriter, r *http.Request) {
	image, err := elasticsearch.GetRandomImage("xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the xkcd source")
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getXkcdByNumber serves a comic by it's number
func getXkcdByNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	image, err := elasticsearch.GetImageByID(vars["number"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getSearchXKCD serves a comic by searching for text
func getSearchXKCD(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image, err := elasticsearch.SearchImage(r.Form["query"], "xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// postXkcdMattermost handles xkcd webhook POST
func postXkcdMattermost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !strings.EqualFold(Token, r.Form["token"][0]) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, errors.New("unauthorized - bad token"))
		log.Error(errors.New("unauthorized - bad token"))
		return
	}

	image, err := elasticsearch.SearchImage(r.Form["trigger_word"], "xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		Text: fmt.Sprintf("### %s\n>%s\n\n%s",
			image.Title,
			image.Text,
			makeURL("http", Host, Port, image.Path),
		),
		Username: "xkcd",
		IconURL:  makeURL("http", Host, Port, "icon/xkcd"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// postXkcdMattermostSlash handles xkcd webhook POST for use with a slash command
func postXkcdMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var image elasticsearch.ImageMetaData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error(errors.New("unauthorized - bad token"))
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random xkcd")
		image, err = elasticsearch.GetRandomImage("xkcd")
	} else if isNumeric(textArg) {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting xkcd by number")
		image, err = elasticsearch.GetImageByID(textArg)
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting xkcd by title")
		image, err = elasticsearch.SearchImage(r.Form["text"], "xkcd")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		ResponseType: "in_channel",
		Text: fmt.Sprintf("### %s\n\non behalf of @%s %s",
			image.Title,
			// image.Text,
			userName,
			makeURL("http", Host, Port, image.Path),
		),
		Username: "xkcd",
		IconURL:  makeURL("http", Host, Port, "icon/xkcd"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// getRandomDilbert serves a random dilbert comic
func getRandomDilbert(w http.ResponseWriter, r *http.Request) {
	image, err := elasticsearch.GetRandomImage("dilbert")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the xkcd source")
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getDilbertByDate serves a comic by it's date
func getDilbertByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	image, err := elasticsearch.GetImageByID(vars["date"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getSearchDilbert serves a comic by searching for text
func getSearchDilbert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image, err := elasticsearch.SearchImage(r.Form["query"], "dilbert")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// postDilbertMattermost handles dilbert webhook POST
func postDilbertMattermost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !strings.EqualFold(Token, r.Form["token"][0]) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, errors.New("unauthorized - bad token"))
		log.Error(errors.New("unauthorized - bad token"))
		return
	}

	image, err := elasticsearch.SearchImage(r.Form["trigger_word"], "dilbert")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		Text: fmt.Sprintf("### %s\n>%s\n\n%s",
			image.Title,
			image.Text,
			makeURL("http", Host, Port, image.Path),
		),
		Username: "dilbert",
		IconURL:  makeURL("http", Host, Port, "icon/dilbert"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// postDilbertMattermostSlash handles dilbert webhook POST for use with a slash command
func postDilbertMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var image elasticsearch.ImageMetaData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error(errors.New("unauthorized - bad token"))
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random dilbert")
		image, err = elasticsearch.GetRandomImage("dilbert")
	} else if isNumeric(textArg) {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting dilbert by number")
		image, err = elasticsearch.GetImageByID(textArg)
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting dilbert by title")
		image, err = elasticsearch.SearchImage(r.Form["text"], "dilbert")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		ResponseType: "in_channel",
		Text: fmt.Sprintf("### %s\n>%s\n\non behalf of @%s %s",
			image.Title,
			image.Text,
			userName,
			makeURL("http", Host, Port, image.Path),
		),
		Username: "dilbert",
		IconURL:  makeURL("http", Host, Port, "icon/dilbert"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// getRandomGiphy serves a random giphy gif
func getRandomGiphy(w http.ResponseWriter, r *http.Request) {
	image, err := elasticsearch.GetRandomImage("giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the giphy source")
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getSearchGiphy serves a comic by searching for text
func getSearchGiphy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image, err := elasticsearch.SearchImage(r.Form["query"], "giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// postGiphyMattermost handles Giphy webhook POST
func postGiphyMattermost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !strings.EqualFold(Token, r.Form["token"][0]) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, errors.New("unauthorized - bad token"))
		log.Error(errors.New("unauthorized - bad token"))
		return
	}

	image, err := elasticsearch.SearchImage(r.Form["trigger_word"], "giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		Text:     makeURL("http", Host, Port, image.Path),
		Username: "scifgif",
		IconURL:  makeURL("http", Host, Port, "icon/giphy"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// postGiphyMattermostSlash handles giphy webhook POST for use with a slash command
func postGiphyMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var image elasticsearch.ImageMetaData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error(errors.New("unauthorized - bad token"))
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random gif")
		image, err = elasticsearch.GetRandomImage("giphy")
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting gif by keyword")
		image, err = elasticsearch.SearchImage(r.Form["text"], "giphy")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		ResponseType: "in_channel",
		Text: fmt.Sprintf("**%s** on behalf of @%s %s",
			textArg,
			userName,
			makeURL("http", Host, Port, image.Path),
		),
		Username: "scifgif",
		IconURL:  makeURL("http", Host, Port, "icon/giphy"),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// getRandomASCII serves a random ascii emoji
func getRandomASCII(w http.ResponseWriter, r *http.Request) {
	ascii, err := elasticsearch.GetRandomASCII()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Error(err)
		return
	}
	log.Debugf("GET %s", ascii.Keywords)
	fmt.Fprint(w, ascii.Emoji)
}

// getSearchASCII serves a comic by searching for text
func getSearchASCII(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ascii, err := elasticsearch.SearchASCII(r.Form["query"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}
	log.Debugf("GET %s", r.Form["query"])
	fmt.Fprint(w, ascii.Emoji)
}

// postASCIIMattermost handles ASCII webhook POST
// func postASCIIMattermost(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
//
// 	if !strings.EqualFold(Token, r.Form["token"][0]) {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintln(w, errors.New("unauthorized - bad token"))
// 		log.Error(errors.New("unauthorized - bad token"))
// 		return
// 	}
//
// 	image, err := elasticsearch.SearchImage(r.Form["trigger_word"], "ascii")
// 	if err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprintln(w, err.Error())
// 		log.Error(err)
// 		return
// 	}
//
// 	webhook := WebHookResponse{
// 		Text:     makeURL("http", Host, Port, image.Path),
// 		Username: "scifgif",
// 		IconURL:  makeURL("http", Host, Port, "icon/giphy"),
// 	}
//
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)
//
// 	if err := json.NewEncoder(w).Encode(webhook); err != nil {
// 		log.Error(err)
// 	}
// }

// postASCIIMattermostSlash handles ascii webhook POST for use with a slash command
func postASCIIMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var ascii elasticsearch.ASCIIData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error(errors.New("unauthorized - bad token"))
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random ascii")
		ascii, err = elasticsearch.GetRandomASCII()
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting ascii by keyword")
		ascii, err = elasticsearch.SearchASCII(r.Form["text"])
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.Error(err)
		return
	}

	webhook := WebHookResponse{
		ResponseType: "in_channel",
		Text:         "# " + ascii.Emoji,
		Username:     userName,
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webhook); err != nil {
		log.Error(err)
	}
}

// getGiphyIcon serves giphy icon
func getGiphyIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "images/icons/giphy-icon.png")
}

// getXkcdIcon serves xkcd icon
func getXkcdIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "images/icons/xkcd-icon.jpg")
}

// getDilbertIcon serves xkcd icon
func getDilbertIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "images/icons/dilbert-icon.png")
}

// getDefaultImage gets default image path
func getDefaultImage(source string) string {
	switch source {
	case "xkcd":
		return "images/default/xkcd.png"
	case "dilbert":
		return "images/default/dilbert.png"
	case "giphy":
		return "images/default/giphy.gif"
	default:
		return "images/default/giphy.gif"
	}
}

// addImage add scifgif GIF
func addImage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
		return
	}

	r.ParseForm()
	keywords := r.Form["keywords"]
	keywords = strings.Split(keywords[0], ",")
	log.Debugf("keywords: %#v", keywords)

	if len(keywords) == 0 {
		http.Error(w, "you forgot to include the keywords parameter", http.StatusBadRequest)
		log.Errorf("no keywords: %v", keywords)
		return
	}
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	defer file.Close()

	log.Debugf("PUT %s/%s", contribFolder, handler.Filename)
	// fmt.Fprintf(w, "%v", handler.Header)

	// protect against file inclusion
	fMD5, err := getMD5(file)
	if err != nil {
		return
	}
	// protect against directory traversal
	fName := filepath.Clean(filepath.Base(fMD5 + filepath.Ext(handler.Filename)))
	fPath := filepath.Join(contribFolder, fName)
	if _, err = os.Stat(fPath); !os.IsNotExist(err) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "file already exists")
		log.Errorf("file already exists: %s", handler.Filename)
		return
	}

	f, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	defer f.Close()

	file.Seek(0, io.SeekStart)
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	elasticsearch.WriteImageToDatabase(elasticsearch.ImageMetaData{
		ID:     fMD5,
		Source: "giphy",
		Name:   filepath.Base(handler.Filename),
		// Title:   ,
		Text: strings.Join(keywords, " "),
		Path: fPath,
	}, "giphy")
	fmt.Fprintln(w, "image successfully added")
}

// updateImageKeywords add keywords to GIF
func updateImageKeywords(w http.ResponseWriter, r *http.Request) {

	if r.Method != "PATCH" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
		return
	}

	vars := mux.Vars(r)
	folder := vars["source"]
	file := vars["file"]

	r.ParseForm()
	keywords := strings.Split(r.FormValue("keywords"), ",")
	log.Debugf("keywords: %#v", keywords)

	// get image by path
	file = filepath.Clean(filepath.Base(file))
	path := filepath.Join("images", folder, file)
	image, err := elasticsearch.GetImageByPath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Error(err)
		return
	}
	// update image's keywords
	image.Text += " " + strings.Join(keywords, " ")
	err = elasticsearch.UpdateKeywords(image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	fmt.Fprintln(w, "image successfully updated")
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["source"]
	file := vars["file"]

	// protect against directory traversal
	file = filepath.Clean(filepath.Base(file))
	path := filepath.Join("images", folder, file)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "image not found")
		log.Error(err)
		return
	}

	log.Infof("deleting images/%s/%s", folder, file)
	err := os.Remove(path)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		log.Error(err, "unable to remove image")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "image successfully removed")
		log.Error(err)
		return
	}
}

// getImage serves scifgif icon
func getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["source"]
	file := vars["file"]

	// protect against directory traversal
	file = filepath.Clean(filepath.Base(file))

	if _, err := os.Stat(filepath.Join("images", folder, file)); os.IsNotExist(err) {
		log.Debugf("GET default %s image", folder)
		http.ServeFile(w, r, getDefaultImage(folder))
		return
	}
	log.Debugf("GET images/%s/%s", folder, file)
	http.ServeFile(w, r, filepath.Join("images", folder, file))
}
