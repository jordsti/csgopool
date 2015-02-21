package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	"csgodb"
)

type ViewEventPage struct {
	Page
	EventName string
	Matches template.HTML
}

func ViewEventHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	t, err := MakeTemplate("viewevent.html")
	if err != nil {
		fmt.Println(err)
	}
	
	_evt_id, _ := strconv.ParseInt(r.FormValue("id"), 10, 32)
	evt_id := int(_evt_id)
	
	db, _ := csgodb.Db.Open()
	
	event := csgodb.GetEventById(db, evt_id)
	event.Matches = csgodb.GetMatchesByEventId(db, event.EventId)
	
	//nil checkup todo	
	matches_html := ""
	
	for _, m := range event.Matches {
		
		t1 := csgodb.GetTeamById(db, m.Team1.TeamId)
		t2 := csgodb.GetTeamById(db, m.Team2.TeamId)
		
		dateStr := fmt.Sprintf("%d-%02d-%02d", m.Date.Year(), m.Date.Month(), m.Date.Day())
		
		mLink := &Link{Caption: dateStr, Url: "/viewmatch/"}
		mLink.AddParameter("id", strconv.Itoa(m.MatchId))
		
		t1cap := fmt.Sprintf("%s (%d)", t1.Name, m.Team1.Score)
		t2cap := fmt.Sprintf("%s (%d)", t2.Name, m.Team2.Score)
		
		t1Link := &Link{Caption: t1cap, Url:"/viewteam/"}
		t1Link.AddParameter("id", strconv.Itoa(t1.TeamId))
		
		t2Link := &Link{Caption: t2cap, Url:"/viewteam/"}
		t2Link.AddParameter("id", strconv.Itoa(t2.TeamId))
		
		matches_html = matches_html + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", mLink.GetHTML() , t1Link.GetHTML(), t2Link.GetHTML(), m.Map)
		
	}
	
	db.Close()
	
	p := &ViewEventPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Last Events"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.EventName = event.Name
	p.Matches = template.HTML(matches_html)
	p.GenerateRightSide(session)
	t.Execute(w, p)
}