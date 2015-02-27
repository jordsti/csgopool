package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"csgodb"
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

	db, _ := csgodb.Db.Open()
	players := csgodb.GetAllPlayersWithStat(db)
	players_html := ""
	
	for _, pl := range players {
		playerLink := &Link{Caption: pl.Name, Url:"/viewplayer/"}
		playerLink.AddInt("id", pl.PlayerId)
		
		players_html += fmt.Sprintf(`<tr>
									<td>%s</td>
									<td>%d</td>
									<td>%d</td>
									<td>%.2f</td>
									<td>%d</td>
									</tr>`, 
									playerLink.GetHTML(), 
									pl.Stat.Frags, 
									pl.Stat.Deaths, 
									pl.Stat.AvgKDRatio, 
									pl.Stat.MatchesPlayed)
	}
	
	db.Close()
	
	p := &PlayersPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Players"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Players = template.HTML(players_html)
	p.GenerateRightSide(session)
	if !session.IsLogged() {
		p.AddLogin(session)
	}
	
	t.Execute(w, p)
}