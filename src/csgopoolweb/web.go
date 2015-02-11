package csgopoolweb

import (
	"net/http"
	"csgopool"
)

var rootPath = ""
var state = &WebServerState{}

type WebServerState struct {
	RootPath string
	Data *csgopool.GameData
	Sessions *SessionContainer
}

func NewWebServer(data *csgopool.GameData, path string) *WebServerState {
	rootPath = path
	state.RootPath = path
	state.Data = data
	state.Sessions = &SessionContainer{}
	
	return state
}

type Page struct {
	Title string
}

func GetMenu() Menu {
	m := Menu{}
	
	i := MenuItem{MenuId: 0, LinkName: "Home", Link:"/"}
	m.Items = append(m.Items, i)
	
	i = MenuItem{MenuId: 1, LinkName: "Last Events", Link:"/events/"}
	m.Items = append(m.Items, i)
	
	i = MenuItem{MenuId: 2, LinkName: "Teams", Link:"/teams/"}
	m.Items = append(m.Items, i)
	
	i = MenuItem{MenuId: 3, LinkName: "Ranking", Link:"/ranking/"}
	m.Items = append(m.Items, i)
	
	return m
}

func (w *WebServerState) Serve() {
	rootPath = w.RootPath
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/viewmatch/", ViewMatchHandler)
	http.HandleFunc("/events/", EventsHandler)
	http.HandleFunc("/viewevent/", ViewEventHandler)
	http.HandleFunc("/teams/", TeamsHandler)
	http.HandleFunc("/viewteam/", ViewTeamHandler)
	http.ListenAndServe(":8080", nil)
}