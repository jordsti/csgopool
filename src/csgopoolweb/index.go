package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgodb"
	"strconv"
	"time"
)

type IndexPage struct {
	Page
	LastUpdate string
	ServerTime string
	Content template.HTML
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	t, err := MakeTemplate("index.html")
	if err != nil {
		fmt.Println(err)
	}
	
	m := GetMenu(session)
	
	curevent := ""
	
	db, _ := csgodb.Db.Open()
	
	event := csgodb.GetLastEvent(db)
	matches := csgodb.GetMatchesByEventId(db, event.EventId)
	last_update := csgodb.GetLastUpdate(db)
	content := ""
	if event != nil {
		matches_html := "<ul>"
		
		for _, m := range matches {
			t1 := csgodb.GetTeamById(db, m.Team1.TeamId)
			t2 := csgodb.GetTeamById(db, m.Team2.TeamId)
			
			matches_html = matches_html + fmt.Sprintf("<li><a href=\"/viewmatch/?id=%d\">(%d) %s vs (%d) %s</a></li>", m.MatchId, m.Team1.Score, t1.Name, m.Team2.Score, t2.Name)
		}
		
		matches_html = matches_html + "</ul>"
		
		evtLink := &Link{Caption:"View Event", Url:"/viewevent/"}
		evtLink.AddParameter("id", strconv.Itoa(event.EventId))
		
		curevent = fmt.Sprintf("<strong>%s</strong><br />%s<br />%s", event.Name, evtLink.GetHTML(), matches_html)
	} else {
		curevent = "<em>No event found !</em>"
	}
	
	last_matches := csgodb.GetLastMatch(db)

	if last_matches.MatchId != 0 {
		stats := csgodb.GetMatchPoints(db, last_matches.MatchId)
		matchLink := &Link{Caption: fmt.Sprintf("%d-%02d-%02d", last_matches.Date.Year(), last_matches.Date.Month(), last_matches.Date.Day()), Url:"/viewmatch/"}
		matchLink.AddInt("id", last_matches.MatchId)
		content = fmt.Sprintf(`Last matches : %s <br />`, matchLink.GetHTML())
		
		content += `<table class="table table-striped"><thead>
			<tr>
				<th>Player</th>
				<th>Team</th>
				<th>Frags</th>
				<th>Headshots</th>
				<th>K/D Ratio</th>
				<th>Points</th>
			</tr>
		</thead><tbody>`
		
		for _, s := range stats {
			pLink := &Link{Caption: s.PlayerName, Url:"/viewplayer/"}
			pLink.AddInt("id", s.PlayerId)
			tLink := &Link{Caption: s.TeamName, Url:"/viewteam/"}
			tLink.AddInt("id", s.TeamId)
			content += fmt.Sprintf(`<tr>
				<td>%s</td>
				<td>%s</td>
				<td>%d</td>
				<td>%d</td>
				<td>%.2f</td>
				<td>%d</td>
			</tr>`, pLink.GetHTML(), tLink.GetHTML(), s.Frags, s.Headshots, s.KDRatio, s.Points)
			
		}
		
		content += `</tbody></table>`
		
		
	} else { 
		content = `<h4>No match found!</h4>`
	}
	
	db.Close()
	
	p := &IndexPage{}
	p.Title = "CS:GO Pool - Home"
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	p.LeftSide = template.HTML(curevent)
	p.LastUpdate = fmt.Sprintf("%02d:%02d", last_update.Time.Hour(), last_update.Time.Minute())
	p.Content = template.HTML(content)
	servertime := time.Now()
	p.ServerTime = fmt.Sprintf("%02d:%02d", servertime.Hour(), servertime.Minute())
	
	p.GenerateRightSide(session)
	
	t.Execute(w, p)
	
}