package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgoscrapper"
	"strconv"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	t, err := MakeTemplate("index.html")
	if err != nil {
		fmt.Println(err)
	}
	
	m := GetMenu(session)
	
	curevent := ""
	
	event := csgoscrapper.GetLastEvent(state.Data.Events)
	
	if event != nil {
		matches := "<ul>"
		
		for _, m := range event.Matches {
			t1 := state.GetTeamById(m.Team1.TeamId)
			t2 := state.GetTeamById(m.Team2.TeamId)
			
			matches = matches + fmt.Sprintf("<li><a href=\"/viewmatch/?id=%d\">(%d) %s vs (%d) %s</a></li>", m.MatchId, m.Team1.Score, t1.Name, m.Team2.Score, t2.Name)
		}
		
		matches = matches + "</ul>"
		
		evtLink := &Link{Caption:"View Event", Url:"/viewevent/"}
		evtLink.AddParameter("id", strconv.Itoa(event.EventId))
		
		curevent = fmt.Sprintf("<strong>%s</strong><br />%s<br />%s", event.Name, evtLink.GetHTML(), matches)
	} else {
		curevent = "<em>No event found !</em>"
	}
	
	p := &Page{}
	p.Title = "CS:GO Pool Home"
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	p.LeftSide = template.HTML(curevent)
	
	if !session.IsLogged() {
		p.AddLogin(session)
	} else {
		p.RightSide = GetUserMenu()
	}
	
	t.Execute(w, p)
	
}