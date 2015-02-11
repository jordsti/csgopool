package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	//"csgoscrapper"
)

type TeamPage struct {
	Title string
	Brand string
	Event string
	Menu template.HTML
	TeamName string
	Players template.HTML
}

func ViewTeamHandler(w http.ResponseWriter, r *http.Request) {
	
	t, err := template.ParseFiles(rootPath + "viewteam.html")
	if err != nil {
		fmt.Println(err)
	}
	
	
	
	
	_m_id := r.FormValue("id")
	m_id, _ := strconv.ParseInt(_m_id, 10, 32)
	
	teamId := int(m_id)
	
	team := state.GetTeamById(teamId)
	
	
	pStats := ""
	for _, pl := range team.Players {
		playerLink := &Link{Caption: pl.Name, Url:"/viewplayer/"}
		playerLink.AddParameter("id", strconv.Itoa(pl.PlayerId))
		
		pStats = pStats + fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%0.2f</td><td>%d</td><td>%.2f</td><td>%d</td><td>%d</td><td>%.2f</td><td>%.2f</td><td>%.2f</td></tr>", playerLink.GetHTML(), pl.Stats.Frags, pl.Stats.Headshot, pl.Stats.Deaths, pl.Stats.KDRatio, pl.Stats.MapsPlayed, pl.Stats.RoundsPlayed, pl.Stats.AvgFragsPerRound, pl.Stats.AvgAssistsPerRound, pl.Stats.AvgDeathsPerRound)
		
	}
	
	
	p := &TeamPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = fmt.Sprintf("CS:GO Pool - Team : %s", team.Name)
	p.Players = template.HTML(pStats)
	p.TeamName = team.Name
	p.Menu = template.HTML(GetMenu().GetHTML())
	
	t.Execute(w, p)


}


	
	
	