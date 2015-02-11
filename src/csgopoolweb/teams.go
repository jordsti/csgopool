package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
)

type TeamsPage struct {
	Title string
	Brand string
	Teams template.HTML
	Menu template.HTML
	
}

func TeamsHandler(w http.ResponseWriter, r *http.Request) {
	
	t, err := template.ParseFiles(rootPath + "teams.html")
	if err != nil {
		fmt.Println(err)
	}

	teams_html := ""
	
	for _, t := range state.Data.Teams {
		
		teamLink := &Link{Caption: t.Name, Url: "/viewteam/"}
		teamLink.AddParameter("id", strconv.Itoa(t.TeamId))
		
		teams_html = teams_html + fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%.2f</td></tr>", teamLink.GetHTML(), t.Stats.Wins, t.Stats.Draws, t.Stats.Losses, t.Stats.Frags, t.Stats.Deaths, t.Stats.RoundsPlayed, t.Stats.KDRatio)
		
	}

	p := &TeamsPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Teams"
	p.Menu = template.HTML(GetMenu().GetHTML())
	p.Teams = template.HTML(teams_html)
	
	t.Execute(w, p)
}