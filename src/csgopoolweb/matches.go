package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"strconv"
	"csgodb"
)

type MatchesPage struct {
	Page
	Matches template.HTML
}

func MatchesHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	t, err := MakeTemplate("matches.html")
	if err != nil {
		fmt.Println(err)
	}
	
	
	db, _ := csgodb.Db.Open()

	matches := csgodb.GetAllMatches(db)
	
	//nil checkup todo	
	matches_html := ""
	
	for _, m := range matches {
		
		
		dateStr := fmt.Sprintf("%d-%02d-%02d", m.Date.Year(), m.Date.Month(), m.Date.Day())
		
		mLink := &Link{Caption: dateStr, Url: "/viewmatch/"}
		mLink.AddParameter("id", strconv.Itoa(m.MatchId))
		
		t1cap := fmt.Sprintf("%s (%d)", m.Team1.Name, m.Team1.Score)
		t2cap := fmt.Sprintf("%s (%d)", m.Team2.Name, m.Team2.Score)
		
		t1Link := &Link{Caption: t1cap, Url:"/viewteam/"}
		t1Link.AddParameter("id", strconv.Itoa(m.Team1.TeamId))
		
		t2Link := &Link{Caption: t2cap, Url:"/viewteam/"}
		t2Link.AddParameter("id", strconv.Itoa(m.Team2.TeamId))
		
		pooled := ""
		
		if m.PoolStatus == 1 {
			pooled = "x"
		}
		
		matches_html = matches_html + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", mLink.GetHTML() , t1Link.GetHTML(), t2Link.GetHTML(), m.Map, pooled, GetMatchLink(m))
		
	}
	
	db.Close()
	
	p := &MatchesPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Last Matches"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Matches = template.HTML(matches_html)
	p.GenerateRightSide(session)
	t.Execute(w, p)
}