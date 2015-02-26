package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgodb"
	"time"
)

type IndexPage struct {
	Page
	LastUpdate string
	ServerTime string
	LastMatch template.HTML
	Divisions template.HTML
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
	
	matches := csgodb.GetMatchesByDate(db, time.Now().AddDate(0, 0, -2))
	last_update := csgodb.GetLastUpdate(db)
	content := ""
	if len(matches) > 0 {
		matches_html := "<ul>"
		
		for _, m := range matches {
			t1 := csgodb.GetTeamById(db, m.Team1.TeamId)
			t2 := csgodb.GetTeamById(db, m.Team2.TeamId)
			
			matches_html = matches_html + fmt.Sprintf("<li><a href=\"/viewmatch/?id=%d\">(%d) %s vs (%d) %s</a></li>", m.MatchId, m.Team1.Score, t1.Name, m.Team2.Score, t2.Name)
		}
		
		matches_html = matches_html + "</ul>"
		
		curevent = fmt.Sprintf("<strong>%s</strong><br />%s", "Matches of the last days", matches_html)
	} else {
		curevent = "<em>No matches found !</em>"
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
				<td>%.2f</td>
				<td>%d</td>
			</tr>`, pLink.GetHTML(), tLink.GetHTML(), s.Frags, s.KDRatio, s.Points)
			
		}
		
		content += `</tbody></table>`
		
		
	} else { 
		content = `<h4>No match found!</h4>`
	}
	
	divisions := csgodb.GetDivisionsPoints(db)
	db.Close()
	/*div_html := ""
	row_html := `<div class="row">%s</div>`
	nb_div := 1
	inner_html := ""
	
	for _, div := range divisions {
		

		inner_html += fmt.Sprintf(`<div class="col-md-4"><h4>%s</h4>`, div.Name)
		
		for _, p := range div.Players {
		
			playerLink := &Link{Caption: p.Name, Url:"/viewplayer/"}
			playerLink.AddInt("id", p.PlayerId)
			
			inner_html += fmt.Sprintf(`%s : %d<br />`, playerLink.GetHTML(), p.Points)
		}
		
		inner_html += `</div>`
		
		if nb_div % 3 == 0 && nb_div > 0 {
			div_html += fmt.Sprintf(row_html, inner_html)
			inner_html = ""
		}
		
		nb_div++
		
	}
	
	if len(inner_html) > 0 {
		div_html += fmt.Sprintf(row_html, inner_html)
	}*/
	//transition with built-in HtmlElement
	
	html_div := CreateDiv()
	html_div.SetAttribute("class", "container")
	nb_div := 0
	currentRow := CreateDiv()
	currentRow.SetAttribute("class", "row")
	html_div.AddChild(currentRow)
	
	for _, div := range divisions {
		
		inner_div := CreateDiv()
		inner_div.SetAttribute("class", "col-md-2")
		
		title := &HtmlElement{Tag: "h4"}
		title.InnerText = div.Name
		
		inner_div.AddChild(title)
		
		for _, p := range div.Players {
			playerLink := &Link{Caption: p.Name, Url:"/viewplayer/"}
			playerLink.AddInt("id", p.PlayerId)
			inner_div.InnerText += fmt.Sprintf(`%s : %d<br />`, playerLink.GetHTML(), p.Points)
		}
		
		if nb_div % 3 == 0 && nb_div > 0 {
			currentRow = CreateDiv()
			currentRow.SetAttribute("class", "row")
			html_div.AddChild(currentRow)
		}
		
		currentRow.AddChild(inner_div)
		
		nb_div++
	}
	
	p := &IndexPage{}
	p.Title = "CS:GO Pool - Home"
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	p.LeftSide = template.HTML(curevent)
	p.LastUpdate = fmt.Sprintf("%02d:%02d", last_update.Time.Hour(), last_update.Time.Minute())
	p.LastMatch = template.HTML(content)
	p.Divisions = template.HTML(html_div.GetHTML())
	servertime := time.Now()
	p.ServerTime = fmt.Sprintf("%02d:%02d", servertime.Hour(), servertime.Minute())
	
	p.GenerateRightSide(session)
	
	t.Execute(w, p)
	
}