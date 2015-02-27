package csgopoolweb

import (
	"net/http"
	"logger"
	"html/template"
	"fmt"
)

var rootPath = ""
var state = &WebServerState{}

type WebServerState struct {
	RootPath string
	Domain string
	Log *logger.LoggerState
	Sessions *SessionContainer
	Port int
}

func NewWebServer(port int, path string, logPath string) *WebServerState {
	rootPath = path
	state.RootPath = path
	state.Sessions = &SessionContainer{}
	state.Domain = "localhost"
	state.Log = &logger.LoggerState{LogPath: logPath, Level: 3}
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

func (p *Page) GenerateRightSide(s *Session) {
	
	if s.IsLogged() {
		p.RightSide = GetUserMenu()
		
	} else {
		if s.IsFieldExists("message") {
			field := s.PopField("message")
			p.Message = template.HTML(field.Value)
		}
		
		p.AddLogin(s)
	}
	
}

func GetMenu(s *Session) Menu {
	m := Menu{}
	
	i := MenuItem{MenuId: 0, LinkName: "Home", Link:"/"}
	m.Items = append(m.Items, i)
	
	i = MenuItem{MenuId: 1, LinkName: "Last Matches", Link:"/matches/"}
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
	http.HandleFunc("/matches/", MatchesHandler)
	http.HandleFunc("/viewuser/", ViewUserHandler)
	http.HandleFunc("/teams/", TeamsHandler)
	http.HandleFunc("/players/", PlayersHandler)
	http.HandleFunc("/viewteam/", ViewTeamHandler)
	http.HandleFunc("/accountform/", AccountFormHandler)
	http.HandleFunc("/createaccount/", CreateAccountHandler)
	http.HandleFunc("/login/", LoginHandler)
	http.HandleFunc("/logout/", LogoutHandler)
	http.HandleFunc("/userpool/", UserPoolHandler)
	http.HandleFunc("/createpool/", CreatePoolHandler)
	http.HandleFunc("/ranking/", RankingHandler)
	
	
	//image serving
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(rootPath + "/images/"))))

	
	http.ListenAndServe(fmt.Sprintf(":%d", w.Port), nil)
}