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

	"github.com/apex/log"
	"github.com/blacktop/scifgif/database"
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

// getWebSearch search for all images matching query
func getWebSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if len(r.Form["query"]) == 0 || len(r.Form["type"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New("bad request - please supply `query` and `type` params")
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("getWebSearch failed")
		return
	}

	images, err := db.SearchGetAll(r.Form["query"], r.Form["type"][0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("getWebSearch failed")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// TODO: only do this for debug mode
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(images); err != nil {
		log.WithError(err).Error("json encoder failed")
	}
}

// getRandomXKCD serves a random xkcd comic
func getRandomXKCD(w http.ResponseWriter, r *http.Request) {

	image, err := db.GetRandomImage("xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the xkcd source")
		log.WithError(err).Error("getRandomXKCD failed")
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getXkcdByNumber serves a comic by it's number
func getXkcdByNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	image, err := db.GetImageByID(vars["number"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("getXkcdByNumber failed")
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getSearchXKCD serves a comic by searching for text
func getSearchXKCD(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image, err := db.SearchImage(r.Form["query"], "xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("getSearchXKCD failed")
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
		log.Error("unauthorized - bad token")
		return
	}

	image, err := db.SearchImage(r.Form["trigger_word"], "xkcd")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("postXkcdMattermost failed")
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
		log.WithError(err).Error("json encoder failed")
	}
}

// postXkcdMattermostSlash handles xkcd webhook POST for use with a slash command
func postXkcdMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var image database.ImageMetaData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error("unauthorized - bad token")
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random xkcd")
		image, err = db.GetRandomImage("xkcd")
	} else if isNumeric(textArg) {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting xkcd by number")
		image, err = db.GetImageByID(textArg)
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting xkcd by title")
		image, err = db.SearchImage(r.Form["text"], "xkcd")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("postXkcdMattermostSlash failed")
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
		log.WithError(err).Error("json encoder failed")
	}
}

// getRandomDilbert serves a random dilbert comic
func getRandomDilbert(w http.ResponseWriter, r *http.Request) {
	image, err := db.GetRandomImage("dilbert")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the xkcd source")
		log.WithError(err).Error("getRandomDilbert failed")
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getDilbertByDate serves a comic by it's date
func getDilbertByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	image, err := db.GetImageByID(vars["date"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("getDilbertByDate failed")
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getSearchDilbert serves a comic by searching for text
func getSearchDilbert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image, err := db.SearchImage(r.Form["query"], "dilbert")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("getSearchDilbert failed")
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
		log.Error("unauthorized - bad token")
		return
	}

	image, err := db.SearchImage(r.Form["trigger_word"], "dilbert")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("postDilbertMattermost failed")
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
		log.WithError(err).Error("json encoder failed")
	}
}

// postDilbertMattermostSlash handles dilbert webhook POST for use with a slash command
func postDilbertMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var image database.ImageMetaData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error("unauthorized - bad token")
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random dilbert")
		image, err = db.GetRandomImage("dilbert")
	} else if isNumeric(textArg) {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting dilbert by number")
		image, err = db.GetImageByID(textArg)
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting dilbert by title")
		image, err = db.SearchImage(r.Form["text"], "dilbert")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("postDilbertMattermostSlash failed")
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
		log.WithError(err).Error("json encoder failed")
	}
}

// getRandomGiphy serves a random giphy gif
func getRandomGiphy(w http.ResponseWriter, r *http.Request) {
	image, err := db.GetRandomImage("giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no images found for the giphy source")
		log.WithError(err).Error("getRandomGiphy failed")
		return
	}
	log.Debugf("GET %s", image.Path)
	http.ServeFile(w, r, image.Path)
}

// getSearchGiphy serves a comic by searching for text
func getSearchGiphy(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	image, err := db.SearchImage(r.Form["query"], "giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("getSearchGiphy failed")
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
		log.Error("unauthorized - bad token")
		return
	}

	image, err := db.SearchImage(r.Form["trigger_word"], "giphy")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("postGiphyMattermost failed")
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
		log.WithError(err).Error("json encoder failed")
	}
}

// postGiphyMattermostSlash handles giphy webhook POST for use with a slash command
func postGiphyMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var image database.ImageMetaData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error("unauthorized - bad token")
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random gif")
		image, err = db.GetRandomImage("giphy")
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting gif by keyword")
		image, err = db.SearchImage(r.Form["text"], "giphy")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("postGiphyMattermostSlash failed")
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
		log.WithError(err).Error("json encoder failed")
	}
}

// getRandomASCII serves a random ascii emoji
func getRandomASCII(w http.ResponseWriter, r *http.Request) {
	ascii, err := db.GetRandomASCII()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.WithError(err).Error("getRandomASCII failed")
		return
	}
	log.Debugf("GET %s", ascii.Keywords)
	fmt.Fprint(w, ascii.Emoji)
}

// getSearchASCII serves a comic by searching for text
func getSearchASCII(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ascii, err := db.SearchASCII(r.Form["query"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("getSearchASCII failed")
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
// 		log.Error("unauthorized - bad token")
// 		return
// 	}
//
// 	image, err := db.SearchImage(r.Form["trigger_word"], "ascii")
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
// 		log.WithError(err).Error("json encoder failed")
// 	}
// }

// postASCIIMattermostSlash handles ascii webhook POST for use with a slash command
func postASCIIMattermostSlash(w http.ResponseWriter, r *http.Request) {
	var ascii database.ASCIIData
	var err error

	r.ParseForm()

	// TODO: add token auth back in

	// if !strings.EqualFold(Token, r.Form["token"][0]) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, errors.New("unauthorized - bad token"))
	// 	log.Error("unauthorized - bad token")
	// 	return
	// }

	userName := r.Form["user_name"][0]
	// teamDomain := r.Form["team_domain"]
	// channelName := r.Form["channel_name"]
	textArg := strings.Join(r.Form["text"], " ")

	if strings.EqualFold(strings.TrimSpace(textArg), "?") {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting random ascii")
		ascii, err = db.GetRandomASCII()
	} else {
		log.WithFields(log.Fields{"text": textArg}).Debug("getting ascii by keyword")
		ascii, err = db.SearchASCII(r.Form["text"])
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("postASCIIMattermostSlash failed")
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
		log.WithError(err).Error("json encoder failed")
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

	// open database
	db, err := database.Open()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "failed to open database")
		log.WithError(err).Error("addImage failed")
		return
	}
	defer db.Close()

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
		log.WithError(err).Error("addImage failed")
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
		log.WithError(err).Error("addImage failed")
		return
	}
	defer f.Close()

	file.Seek(0, io.SeekStart)
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.WithError(err).Error("addImage failed")
		return
	}

	db.WriteImageToDatabase(database.ImageMetaData{
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

	if len(keywords) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "must supply keywords")
		return
	}

	// get image by path
	file = filepath.Clean(filepath.Base(file))
	path := filepath.Join("images", folder, file)
	image, err := db.GetImageByPath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.WithError(err).Error("updateImageKeywords failed")
		return
	}
	// update image's keywords
	image.Text += " " + strings.Join(keywords, " ")
	err = db.UpdateKeywords(image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.WithError(err).Error("updateImageKeywords failed")
		return
	}
	fmt.Fprintln(w, "image successfully updated")
}

func deleteImage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
		return
	}

	vars := mux.Vars(r)
	folder := vars["source"]
	file := vars["file"]

	// protect against directory traversal
	file = filepath.Clean(filepath.Base(file))
	path := filepath.Join("images", folder, file)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "image not found")
		log.WithError(err).Error("deleteImage failed")
		return
	}

	log.Infof("deleting images/%s/%s", folder, file)
	err := os.Remove(path)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		log.WithError(err).Error("unable to remove image")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "image successfully removed")
		log.WithError(err).Error("deleteImage")
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
