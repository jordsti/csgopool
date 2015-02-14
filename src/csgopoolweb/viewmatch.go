package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	"csgoscrapper"
)

type MatchPage struct {
	Page
	Menu template.HTML
	Map string
	PlayerStats template.HTML
}

func ViewMatchHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	t, err := MakeTemplate("viewmatch.html")
	if err != nil {
		fmt.Println(err)
	}
	
	
	
	
	_m_id := r.FormValue("id")
	m_id, _ := strconv.ParseInt(_m_id, 10, 32)
	
	matchId := int(m_id)
	
	match := state.GetMatchById(matchId)
	
	t1 := state.GetTeamById(match.Team1.TeamId)
	t2 := state.GetTeamById(match.Team2.TeamId)
	
	if match == nil {
		state.Log.Error(fmt.Sprintf("Match [%d] not found", match.MatchId))
	}
	
	//generating stats
	pStats := `<table class="table table-striped"><thead><tr><th>Player</th><th>Team</th><th>Frags</th><th>Headshots</th><th>Assists</th><th>Deaths</th><th>K/D</th><th>K/D Delta</th></tr></thead><tbody>`
	for _, ps := range match.PlayerStats {
		pl := &csgoscrapper.Player{}
		team := t1
		if ps.TeamId == t1.TeamId {
			pl = t1.GetPlayerById(ps.PlayerId)
		} else {
			pl = t2.GetPlayerById(ps.PlayerId)
			team = t2
		}
		
		teamLink := &Link{Caption: team.Name, Url:"/viewteam/"}
		teamLink.AddParameter("id", strconv.Itoa(team.TeamId))
		
		if pl != nil {
			playerLink := &Link{Caption:pl.Name, Url:"/viewplayer/"}
			playerLink.AddParameter("id", strconv.Itoa(pl.PlayerId))
			playerLink.AddParameter("teamid", strconv.Itoa(team.TeamId))
			
			pStats = pStats + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td><td>%.2f</td><td>%d</td></tr>", playerLink.GetHTML(), teamLink.GetHTML(), ps.Frags, ps.Headshots, ps.Assists, ps.Deaths, ps.KDRatio, ps.KDDelta)
		}
	}
	
	pStats = pStats + "</tbody></table>"
	
	p := &MatchPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = fmt.Sprintf("CS:GO Pool - Match : %s versus %s", t1.Name, t2.Name)
	p.PlayerStats = template.HTML(pStats)
	p.Map = match.Map
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	
	t.Execute(w, p)


}


	
	
	