package csgopool

import (
	"fmt"
	"csgoscrapper"
)

func InitTeams(path string) []*csgoscrapper.Team {
	
	teams := csgoscrapper.LoadTeams(path)
	
	if len(teams) == 0 {
		fmt.Println("Team not found! Fetching teams from HLTV.org")
		
		p := csgoscrapper.GetTeamsPage()
		
		content, _ := p.LoadPage()
		
		teams = content.ParseTeams()
		
		for _, t := range teams {
			t.LoadTeam()
		}
		
		fmt.Printf("%d teams loaded !", len(teams))
		
		csgoscrapper.SaveTeams(teams, path)
	}
	
	return teams
	
}  

func InitEvents(path string) []*csgoscrapper.Event {
	events := csgoscrapper.LoadEvents(path)
	
	if len(events) == 0 {
		fmt.Println("Events not found! Fetching events from HLTV.org")
		
		page := csgoscrapper.GetEventsPage()
	
		pc, _ := page.LoadPage()
		
		evts := pc.ParseEvents()
		
		for _, e := range evts {
			fmt.Printf("%d, %s\n", e.EventId, e.Name)
			e.LoadAllMatches()
		}
	}
	
	csgoscrapper.SaveEvents(events, path)
	
	return events
	
}