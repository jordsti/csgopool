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
	Matches template.HTML
	Source template.HTML
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
	state.Log.Debug(fmt.Sprintf("Players Count : %d", len(team.Players)))
	for _, pl := range team.Players {
		playerLink := &Link{Caption: pl.Name, Url:"/viewplayer/"}
		playerLink.AddInt("id", pl.PlayerId)
		//playerLink.AddInt("teamid", team.TeamId)
		
		pStats += fmt.Sprintf(`<tr>
									<td>%s</td>
									<td>%d</td>
									<td>%d</td>
									<td>%.2f</td>
									<td>%d</td>
									<td>%.2f</td>
									<td>%.2f</td>
								</tr>`, 
								playerLink.GetHTML(), 
								pl.Stat.Frags, 
								pl.Stat.Deaths, 
								pl.Stat.AvgKDRatio, 
								pl.Stat.MatchesPlayed, 
								pl.Stat.AvgFrags, 
								pl.Stat.AvgKDDelta)
		
	}
	
	matches := csgodb.GetTeamMatches(db, team.TeamId)
	matches_html := ""
	for _, m := range matches {
		matchLink := &Link{Caption: fmt.Sprintf("%d-%02d-%02d", m.Date.Year(), m.Date.Month(), m.Date.Day()), Url:"/viewmatch/"}
		matchLink.AddInt("id", m.MatchId)
		
		team1Link := &Link{Caption: fmt.Sprintf("%s (%d)", m.Team1.Name, m.Team1.Score), Url:"/viewteam/"}
		team1Link.AddInt("id", m.Team1.TeamId)
		
		team2Link := &Link{Caption: fmt.Sprintf("%s (%d)", m.Team2.Name, m.Team2.Score), Url:"/viewteam/"}
		team2Link.AddInt("id", m.Team2.TeamId)
		
		matches_html += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", matchLink.GetHTML(), team1Link.GetHTML(), team2Link.GetHTML(), m.Map)
	}
	
	
	db.Close()
	
	p := &TeamPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = fmt.Sprintf("CS:GO Pool - Team : %s", team.Name)
	p.Players = template.HTML(pStats)
	p.TeamName = team.Name
	p.Source = template.HTML(GetTeamLink(&team_))
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Matches = template.HTML(matches_html)
	p.GenerateRightSide(session)
	t.Execute(w, p)


}


	
	
	