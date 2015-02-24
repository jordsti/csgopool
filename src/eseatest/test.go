package main
import (
	"eseascrapper"
	"time"
	"fmt"
)

func main() {
	
	date := time.Date(2015, 2, 22, 0, 0, 0, 0, time.Local)
	
	url := eseascrapper.GetScheduleURL(date, "invite", eseascrapper.NorthAmerica)
	
	page := url.LoadPage()
	
	
	matches := page.ParseMatches()
	
	for _, m := range matches {
		m.ParseMatch()
		
		fmt.Printf("Match [%d], Team1 [%d] %d, Team2 [%d] %d\n",m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score)
	}
	
}