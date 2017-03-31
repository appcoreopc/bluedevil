package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type MediaInfo struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*MediaInfo, error) {
	filename := title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &MediaInfo{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "video/webm")
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p.Body)

}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
