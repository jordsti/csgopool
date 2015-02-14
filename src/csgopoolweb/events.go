package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
)

type EventsPage struct {
	Page
	Events template.HTML
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	t, err := MakeTemplate("events.html")
	if err != nil {
		fmt.Println(err)
	}
	
	
	evts_html := ""
	
	for _, evt := range state.Data.Events {
		
		evtLink := &Link{Caption: evt.Name, Url: "/viewevent/"}
		evtLink.AddParameter("id", strconv.Itoa(evt.EventId))
		
		evts_html = evts_html + fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", evtLink.GetHTML(), len(evt.Matches))
	}
	
	p := &EventsPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Last Events"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Events = template.HTML(evts_html)
	
	t.Execute(w, p)
}