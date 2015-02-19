package csgopoolweb

import (
	"net/http"
	"csgopool"
	"csgoscrapper"
	"html/template"
	"fmt"
)

var rootPath = ""
var state = &WebServerState{}

type WebServerState struct {
	RootPath string
	Domain string
	Log *csgoscrapper.LoggerState
	Data *csgopool.GameData
	Sessions *SessionContainer
	Users *csgopool.Users
	Port int
}

func NewWebServer(data *csgopool.GameData, users *csgopool.Users, port int, path string, logPath string) *WebServerState {
	rootPath = path
	state.RootPath = path
	state.Data = data
	state.Sessions = &SessionContainer{}
	state.Domain = "localhost"
	state.Users = users
	state.Log = &csgoscrapper.LoggerState{LogPath: logPath, Level: 3}
	state.Port = port
	return state
}

type Page struct {
	Title string
	Brand string
	RightSide template.HTML
	LeftSide template.HTML
	Menu template.HTML
	Message template.HTML
}

func (p *Page) AddLogin(s *Session) {
	if s.IsFieldExists("message") {
		field := s.PopField("message")
		p.Message = template.HTML(field.Value)
	}
	
	p.RightSide = GetLoginForm()
}

func GetMenu(s *Session) Menu {
	m := Menu{}
	
	i := MenuItem{MenuId: 0, LinkName: "Home", Link:"/"}
	m.Items = append(m.Items, i)
	
	i = MenuItem{MenuId: 1, LinkName: "Last Events", Link:"/events/"}
	m.Items = append(m.Items, i)
	
	i = MenuItem{MenuId: 2, LinkName: "Teams", Link:"/teams/"}
	m.Items = append(m.Items, i)
		
	i = MenuItem{MenuId: 3, LinkName: "Players", Link:"/players/"}
	m.Items = append(m.Items, i)
	
	
	i = MenuItem{MenuId: 4, LinkName: "Ranking", Link:"/ranking/"}
	m.Items = append(m.Items, i)
	
	if s.IsLogged() {
		
		if s.User.IsPoolMaster() {
			i = MenuItem{MenuId: 6, LinkName: "Pool Admin", Link:"/adminpool/?action=menu"}
			m.Items = append(m.Items, i)
		}
		
		i = MenuItem{MenuId: 5, LinkName: "Log out", Link:"/logout/?out"}
		m.Items = append(m.Items, i)
	}
	
	return m
}

func (w *WebServerState) Serve() {
	rootPath = w.RootPath
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/adminpool/", AdminPoolHandler)
	http.HandleFunc("/viewmatch/", ViewMatchHandler)
	http.HandleFunc("/viewplayer/", ViewPlayerHandler)
	http.HandleFunc("/events/", EventsHandler)
	http.HandleFunc("/viewevent/", ViewEventHandler)
	http.HandleFunc("/teams/", TeamsHandler)
	http.HandleFunc("/players/", PlayersHandler)
	http.HandleFunc("/viewteam/", ViewTeamHandler)
	http.HandleFunc("/accountform/", AccountFormHandler)
	http.HandleFunc("/createaccount/", CreateAccountHandler)
	http.HandleFunc("/login/", LoginHandler)
	http.HandleFunc("/logout/", LogoutHandler)
	
	
	//image serving
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(rootPath + "/images/"))))

	
	http.ListenAndServe(fmt.Sprintf(":%d", w.Port), nil)
}