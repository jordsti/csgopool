package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"csgodb"
)

type TeamsPage struct {
	Page
	Teams template.HTML
}

func TeamsHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	t, err := MakeTemplate("teams.html")
	if err != nil {
		fmt.Println(err)
	}

	teams_html := ""
	
	db, _ := csgodb.Db.Open()
	
	teams :=  csgodb.GetTeamsWithCount(db)
	db.Close()
	
	for _, t := range teams {
		if t.MatchesCount > 0 {
			teamLink := &Link{Caption: t.Name, Url: "/viewteam/"}
			teamLink.AddInt("id", t.TeamId)
			
			teams_html += fmt.Sprintf(`
									<tr>
										<td>%s</td>
										<td>%d</td>
										<td>%d</td>
									</tr>`, 
									teamLink.GetHTML(), 
									t.PlayersCount, 
									t.MatchesCount)
		}
	}

	p := &TeamsPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Teams"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Teams = template.HTML(teams_html)
	p.GenerateRightSide(session)
	if !session.IsLogged() {
		p.AddLogin(session)
	}
	
	t.Execute(w, p)
}