package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgodb"
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
	
	db, _ := csgodb.Db.Open()
	
	event := csgodb.GetLastEvent(db)
	matches := csgodb.GetMatchesByEventId(db, event.EventId)
	
	if event != nil {
		matches_html := "<ul>"
		
		for _, m := range matches {
			t1 := csgodb.GetTeamById(db, m.Team1.TeamId)
			t2 := csgodb.GetTeamById(db, m.Team2.TeamId)
			
			matches_html = matches_html + fmt.Sprintf("<li><a href=\"/viewmatch/?id=%d\">(%d) %s vs (%d) %s</a></li>", m.MatchId, m.Team1.Score, t1.Name, m.Team2.Score, t2.Name)
		}
		
		matches_html = matches_html + "</ul>"
		
		evtLink := &Link{Caption:"View Event", Url:"/viewevent/"}
		evtLink.AddParameter("id", strconv.Itoa(event.EventId))
		
		curevent = fmt.Sprintf("<strong>%s</strong><br />%s<br />%s", event.Name, evtLink.GetHTML(), matches_html)
	} else {
		curevent = "<em>No event found !</em>"
	}
	
	db.Close()
	
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