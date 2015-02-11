package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
)

type ViewEventPage struct {
	Title string
	Brand string
	EventName string
	Menu template.HTML
	Matches template.HTML
}

func ViewEventHandler(w http.ResponseWriter, r *http.Request) {
	
	t, err := template.ParseFiles(rootPath + "viewevent.html")
	if err != nil {
		fmt.Println(err)
	}
	
	_evt_id, _ := strconv.ParseInt(r.FormValue("id"), 10, 32)
	evt_id := int(_evt_id)
	
	event := state.GetEventById(evt_id)
	//nil checkup todo	
	matches_html := ""
	
	for _, m := range event.Matches {
		
		t1 := state.GetTeamById(m.Team1.TeamId)
		t2 := state.GetTeamById(m.Team2.TeamId)
		
		
		mLink := &Link{Caption: m.Date.String(), Url: "/viewmatch/"}
		mLink.AddParameter("id", strconv.Itoa(m.MatchId))
		
		t1cap := fmt.Sprintf("%s (%d)", t1.Name, m.Team1.Score)
		t2cap := fmt.Sprintf("%s (%d)", t2.Name, m.Team2.Score)
		
		t1Link := &Link{Caption: t1cap, Url:"/viewteam/"}
		t1Link.AddParameter("id", strconv.Itoa(t1.TeamId))
		
		t2Link := &Link{Caption: t2cap, Url:"/viewteam/"}
		t2Link.AddParameter("id", strconv.Itoa(t2.TeamId))
		
		matches_html = matches_html + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", mLink.GetHTML() , t1Link.GetHTML(), t2Link.GetHTML(), m.Map)
		
	}

	p := &ViewEventPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Last Events"
	p.Menu = template.HTML(GetMenu().GetHTML())
	p.EventName = event.Name
	p.Matches = template.HTML(matches_html)
	
	t.Execute(w, p)
}