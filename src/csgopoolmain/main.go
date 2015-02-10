package main

import (
	"fmt"
	"csgopool"
	"csgoscrapper"
	"os"
)

func main() {
	fmt.Println("CS GO Pool")
	
	fmt.Println(os.TempDir())
	
	path := os.TempDir() + "/teams.json"
	
	teams := csgopool.InitTeams(path)
	
	/*for _, t := range teams {
		fmt.Printf("%s, %d\n", t.Name, t.PlayersCount())
		for _, p := range t.Players {
			fmt.Printf("	%s, %d\n", p.Name, p.Stats.Frags)
		}
	}*/
	
	fmt.Printf("Teams count : %d\n", len(teams))
	
	page := csgoscrapper.GetEventsPage()
	
	pc, _ := page.LoadPage()
	
	evts := pc.ParseEvents()
	
	for _, e := range evts {
		fmt.Printf("%d, %s\n", e.EventId, e.Name)
		e.LoadAllMatches()
	}
	
	csgoscrapper.SaveEvents(evts, os.TempDir() + "/events.json")
	
}