package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	//"csgoscrapper"
)

type PlayerPage struct {
	Page
	PlayerName string
	TeamName template.HTML
	Matches template.HTML
	Frags string
	Headshots string
	Deaths string
	KDRatio string
	MapsPlayed string
	RoundsPlayed string
	AvgFragsPerRound string
	AvgAssistsPerRound string
	AvgDeathsPerRound string
}

func ViewPlayerHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	t, err := MakeTemplate("viewplayer.html")
	if err != nil {
		fmt.Println(err)
	}
	
	
	
	
	_m_id := r.FormValue("id")
	_t_id := r.FormValue("teamid")
	m_id, _ := strconv.ParseInt(_m_id, 10, 32)
	t_id, _ := strconv.ParseInt(_t_id, 10, 32)
	
	playerId := int(m_id)
	teamId := int(t_id)
	
	team := state.GetTeamById(teamId)
	player := team.GetPlayerById(playerId)
	
	match_html := ""
	
	for _, evt := range state.Data.Events {
	  
	    evtLink := &Link{Caption: evt.Name, Url:"/viewevent/"}
	    evtLink.AddInt("id", evt.EventId)
	  
	  for _, m := range evt.Matches {
	    
	    if m.IsPlayerIn(player.PlayerId) {
	      //add match
	      t1 := state.GetTeamById(m.Team1.TeamId)
	      t2 := state.GetTeamById(m.Team2.TeamId)
	      
	      matchName := fmt.Sprintf("%s - %s (%d) vs %s (%d)", m.Date.String(), t1.Name, m.Team1.Score, t2.Name, m.Team2.Score)
	      
	      mLink := &Link{Caption: matchName, Url:"/viewmatch/" }
	      mLink.AddParameter("id", strconv.Itoa(m.MatchId))
	      
	      match_html = match_html + fmt.Sprintf("%s %s %s<br />", mLink.GetHTML(), m.Map, evtLink.GetHTML())
	    }
	    
	  }
	  
	}
	
	teamLink := &Link{Caption: team.Name, Url:"/viewteam/"}
	teamLink.AddInt("id", team.TeamId)
	
	p := &PlayerPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = fmt.Sprintf("CS:GO Pool - Team : %s", team.Name)
	p.PlayerName = player.Name
	p.TeamName = template.HTML(teamLink.GetHTML())
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Matches = template.HTML(match_html)
	
	p.Frags = fmt.Sprintf("%d", player.Stats.Frags)
	p.Headshots = fmt.Sprintf("%.2f", player.Stats.Headshots)
	p.Deaths = fmt.Sprintf("%d", player.Stats.Deaths)
	p.KDRatio = fmt.Sprintf("%.2f", player.Stats.KDRatio)
	p.MapsPlayed = fmt.Sprintf("%d", player.Stats.MapsPlayed)
	p.RoundsPlayed = fmt.Sprintf("%d", player.Stats.RoundsPlayed)
	
	p.AvgFragsPerRound = fmt.Sprintf("%.2f", player.Stats.AvgFragsPerRound)
	p.AvgAssistsPerRound = fmt.Sprintf("%.2f", player.Stats.AvgAssistsPerRound)
	p.AvgDeathsPerRound = fmt.Sprintf("%.2f", player.Stats.AvgDeathsPerRound)
	
	t.Execute(w, p)


}


	
	
	