package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	"csgodb"
)

type PlayersPage struct {
	Page
	LinkPages template.HTML
	Players template.HTML
}

func PlayersHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	t, err := MakeTemplate("players.html")
	if err != nil {
		fmt.Println(err)
	}

	db, _ := csgodb.Db.Open()
	
	count := 50
	start := 0
	
	str_start := r.FormValue("start")
	
	if len(str_start) > 0 {
		_start, _ := strconv.ParseInt(str_start, 10, 32)
		start = int(_start)
	}
	
	players := csgodb.GetPlayersWithStat(db, start, count)
	players_html := ""
	
	link_pages := ""
	
	if start > 0 {
		prevLink := &Link{Caption:"Previous", Url:"/players/"}
		prevLink.AddInt("start", start-count)
		link_pages += prevLink.GetHTML()
	}
	
	if len(players) == count {
		nextLink := &Link{Caption:"Next", Url:"/players/"}
		nextLink.AddInt("start", start+count)
		if len(link_pages) > 0 {
			link_pages += " | "
		}
		
		link_pages += nextLink.GetHTML()
	}
	
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
	p.LinkPages = template.HTML(link_pages)
	p.Players = template.HTML(players_html)
	p.GenerateRightSide(session)
	if !session.IsLogged() {
		p.AddLogin(session)
	}
	
	t.Execute(w, p)
}