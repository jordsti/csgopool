package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
)

type PlayersPage struct {
	Page
	Players template.HTML
}

func PlayersHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	t, err := MakeTemplate("players.html")
	if err != nil {
		fmt.Println(err)
	}

	players_html := ""
	
	
	
	
	for _, t := range state.Data.Teams {


		for _, p := range t.Players {
			playerLink := &Link{Caption: p.Name, Url:"/viewplayer/"}
			playerLink.AddInt("id", p.PlayerId)
			playerLink.AddInt("teamid", t.TeamId)

			players_html = players_html + fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%.2f</td><td>%d</td><td>%.2f</td><td>%d</td><td>%d</td><td>%.2f</td><td>%.2f</td><td>%.2f</td></tr>", playerLink.GetHTML(), p.Stats.Frags, p.Stats.Headshots, p.Stats.Deaths, p.Stats.KDRatio, p.Stats.MapsPlayed, p.Stats.RoundsPlayed, p.Stats.AvgFragsPerRound, p.Stats.AvgAssistsPerRound, p.Stats.AvgDeathsPerRound)
		}
	}

	p := &PlayersPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Players"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Players = template.HTML(players_html)
	
	if !session.IsLogged() {
		p.AddLogin(session)
	}
	
	t.Execute(w, p)
}