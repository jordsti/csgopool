package csgopoolweb

import (
	"net/http"
	"html/template"
	"csgopool"
	"fmt"
)

var rootPath = ""
var state = &WebServerState{}

type WebServerState struct {
	
	RootPath string
	Data *csgopool.GameData
}

func NewWebServer(data *csgopool.GameData, path string) *WebServerState {
	rootPath = path
	state.RootPath = path
	state.Data = data
	
	return state
}

type Page struct {
	Title string
}

type IndexPage struct {
	Title string
	LastEvents string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	
	t, err := template.ParseFiles(rootPath + "index.html")
	if err != nil {
		fmt.Println(err)
	}
	p := &IndexPage{Title: "CS Go Pool Home", LastEvents: "test events!"}
	t.Execute(w, p)
	
}

func (w *WebServerState) Serve() {
	rootPath = w.RootPath
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8080", nil)
}