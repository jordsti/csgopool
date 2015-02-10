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