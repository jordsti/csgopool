package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	"csgodb"
	//"csgoscrapper"
)

type TeamPage struct {
	Page
	TeamName string
	Players template.HTML
}

func ViewTeamHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	t, err := MakeTemplate("viewteam.html")
	if err != nil {
		fmt.Println(err)
	}
	
	_m_id := r.FormValue("id")
	m_id, _ := strconv.ParseInt(_m_id, 10, 32)
	
	teamId := int(m_id)
	
	db, _ := csgodb.Db.Open()
	
	team_ := csgodb.GetTeamById(db, teamId)
	if team_.TeamId == 0 {
		//error here!
		return
	}
	
	team := team_.P()
	team.FetchPlayers(db)
	pStats := ""
	for _, pl := range team.Players {
		playerLink := &Link{Caption: pl.Name, Url:"/viewplayer/"}
		playerLink.AddInt("id", pl.PlayerId)
		playerLink.AddInt("teamid", team.TeamId)
		
		pStats = pStats + fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%0.2f</td><td>%d</td><td>%.2f</td><td>%d</td><td>%d</td><td>%.2f</td></tr>", playerLink.GetHTML(), pl.Stat.Frags, pl.Stat.Headshots, pl.Stat.Deaths, pl.Stat.AvgKDRatio, pl.Stat.MatchesPlayed, pl.Stat.AvgFrags, pl.Stat.AvgKDDelta)
		
	}
	
	
	db.Close()
	
	p := &TeamPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = fmt.Sprintf("CS:GO Pool - Team : %s", team.Name)
	p.Players = template.HTML(pStats)
	p.TeamName = team.Name
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	
	t.Execute(w, p)


}


	
	
	