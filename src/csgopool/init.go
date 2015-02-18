package csgopool

import (
	"fmt"
	"csgoscrapper"
)

func InitTeams() []*csgoscrapper.Team {
	
	teams := []*csgoscrapper.Team{}
	
	fmt.Println("Team initial import, fetching teams from HLTV.org")
	
	p := csgoscrapper.GetTeamsPage()
	
	content, _ := p.LoadPage()
	
	teams = content.ParseTeams()
	
	for _, t := range teams {
		t.LoadTeam()
	}
	//todo log
	fmt.Printf("%d teams loaded !\n", len(teams))
		
	return teams
	
}  

func InitEvents() []*csgoscrapper.Event {
	events := []*csgoscrapper.Event{}
	

	fmt.Println("Events initial import, fetching events from HLTV.org")
	
	page := csgoscrapper.GetEventsPage()

	pc, _ := page.LoadPage()
	
	events = pc.ParseEvents()
	
	for _, e := range events {
		fmt.Printf("%d, %s\n", e.EventId, e.Name)
		//e.LoadAllMatches()
	}
	/*} else {
		//updating
		
		fmt.Println("Updating events and matches from HLTV.org")
		
		page := csgoscrapper.GetEventsPage()
		
		pc, _ := page.LoadPage()
		
		events = pc.UpdateEvents(events)
		
	}*/
	
	
	return events
	
}