package csgopoolweb

import (
	"net/http"
	"html/template"
	"csgodb"
	"fmt"
	"time"
)

type RankingPage struct {
	Page
	WeekTop10 template.HTML
	Users template.HTML
	Players template.HTML
}

func RankingHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	t, err := MakeTemplate("ranking.html")
	if err != nil {
		fmt.Println(err)
	}
	
	m := GetMenu(session)
	
	db, _ := csgodb.Db.Open()
	
	users_html := ""
	players_html := ""
	week_html := ""
	
	now := time.Now()
	start := now.AddDate(0, 0, -7)
	
	weeks := csgodb.GetUserPointsInRange(db, start, now, 10)
	pos := 1
	for _, w := range weeks {
		userLink := &Link{Caption: w.Name, Url:"/viewuser/"}
		userLink.AddInt("id", w.UserId)
		week_html += fmt.Sprintf(
			`<tr>
				<td>%d</td>
				<td>%s</td>
				<td>%d</td>
			</tr>`,
			pos,
			userLink.GetHTML(),
			w.Points)
		pos++
		
	}
	
	users := csgodb.GetUserPoints(db)
	pos = 1
	for _, u := range users {
		userLink := &Link{Caption: u.Name, Url:"/viewuser/"}
		userLink.AddInt("id", u.UserId)
		users_html += fmt.Sprintf(
			`<tr>
				<td>%d</td>
				<td>%s</td>
				<td>%d</td>
			</tr>`,
			pos,
			userLink.GetHTML(),
			u.Points)
		pos++
	}
	
	players := csgodb.GetPlayersPoint(db, 0, 25)
	pos = 1
	
	for _, pl := range players {
		playerLink := &Link{Caption: pl.Name, Url:"/viewplayer/"}
		playerLink.AddInt("id", pl.PlayerId)
		
		players_html += fmt.Sprintf(`
		<tr>
			<td>%d</td>
			<td>%s</td>
			<td>%d</td>
			<td>%d</td>
			<td>%.2f</td>
			<td>%.2f</td>
			<td>%d</td>
		</tr>
		`,
		pos, 
		playerLink.GetHTML(),
		pl.Matches,
		pl.Frags,
		pl.KDRatio,
		pl.KDDelta,
		pl.Points)
		
		
		pos++
	}
	
	
	db.Close()
	
	p := &RankingPage{}
	p.Title = "CS:GO Pool - Ranking"
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	p.WeekTop10 = template.HTML(week_html)
	p.Users = template.HTML(users_html)
	p.Players = template.HTML(players_html)
	//p.LeftSide = template.HTML(curevent)
	p.GenerateRightSide(session)
	
	t.Execute(w, p)
	
}