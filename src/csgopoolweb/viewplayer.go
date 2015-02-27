package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	"csgodb"
	//"csgoscrapper"
)

type PlayerPage struct {
	Page
	PlayerName string
	Matches template.HTML
	Frags string
	Source template.HTML
	Deaths string
	KDRatio string
	MatchesPlayed string
	RoundsPlayed string
	AvgFrags string
	AvgKDDelta string
	AvgDeathsPerRound string
	
	TeamsStats template.HTML
}

func ViewPlayerHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	t, err := MakeTemplate("viewplayer.html")
	if err != nil {
		fmt.Println(err)
	}
	
	
	db, _ := csgodb.Db.Open()
	
	_m_id := r.FormValue("id")
	m_id, _ := strconv.ParseInt(_m_id, 10, 32)

	playerId := int(m_id)
	
	player := csgodb.GetPlayerWithStatById(db, playerId)
	
	teams_html := ""
	
	teamsStats := csgodb.GetPlayerTeamStats(db, playerId)
	
	for _, t := range teamsStats {
		
		teamLink := &Link{Caption:t.Name, Url:"/viewteam/"}
		teamLink.AddInt("id", t.TeamId)
		
		teams_html += fmt.Sprintf(`<tr>
								<td>%s</td>
								<td>%d</td>
								</tr>`, 
								teamLink.GetHTML(), 
								t.MatchesCount)
	}
	
	matchStats := csgodb.GetPlayerMatchStats(db, playerId)
	
	matches_html := ""
	
	for _, ms := range matchStats {
		team1Link := &Link{Caption: fmt.Sprintf("%s (%d)", ms.Team1.Name, ms.TeamScore1), Url:"/viewteam/"}
		team1Link.AddInt("id", ms.Team1.TeamId)
		
		team2Link := &Link{Caption: fmt.Sprintf("%s (%d)", ms.Team2.Name, ms.TeamScore2), Url:"/viewteam/"}
		team2Link.AddInt("id", ms.Team2.TeamId)
		
		matchLink := &Link{Caption: fmt.Sprintf("%d-%02d-%02d", ms.Date.Year(), ms.Date.Month(), ms.Date.Day()), Url:"/viewmatch/"}
		matchLink.AddInt("id", ms.MatchId)
		
		matches_html += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%.2f</td><td>%d</td></tr>", matchLink.GetHTML(), team1Link.GetHTML(), team2Link.GetHTML(), ms.Frags, ms.KDRatio, ms.Points)
	}
	
	db.Close()
	
	p := &PlayerPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = fmt.Sprintf("CS:GO Pool - Player : %s", player.Name)
	p.PlayerName = player.Name

	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Matches = template.HTML(matches_html)
	
	
	p.Frags = fmt.Sprintf("%d", player.Stat.Frags)
	//p.Headshots = fmt.Sprintf("%d", player.Stat.Headshots)
	p.Deaths = fmt.Sprintf("%d", player.Stat.Deaths)
	p.KDRatio = fmt.Sprintf("%.2f", player.Stat.AvgKDRatio)
	p.MatchesPlayed = fmt.Sprintf("%d", player.Stat.MatchesPlayed)
	
	p.AvgFrags = fmt.Sprintf("%.2f", player.Stat.AvgFrags)
	p.AvgKDDelta = fmt.Sprintf("%.2f", player.Stat.AvgKDDelta)
	p.Source = template.HTML(GetPlayerLink(&player.Player))
	p.TeamsStats = template.HTML(teams_html)
	p.GenerateRightSide(session)
	t.Execute(w, p)


}


	
	
	