package csgopool

import (
	"fmt"
	"csgoscrapper"
	"time"
)



type GameData struct {
	Events []*csgoscrapper.Event
	Teams []*csgoscrapper.Team
	MissingTeams []int
}

type WatcherState struct {
	
	DataPath string
	Data GameData
	Running bool
	
}



func NewWatcher(dataPath string) *WatcherState {
	
	state := &WatcherState{DataPath: dataPath}
	state.Running = false
	return state
}

func (w *WatcherState) LoadData() {
	
	w.Running = true
	
	team_path := w.DataPath + "/teams.json"
	
	teams := InitTeams(team_path)
	
	w.Data.Teams = teams
	/*for _, t := range teams {
		fmt.Printf("%s, %d\n", t.Name, t.PlayersCount())
		for _, p := range t.Players {
			fmt.Printf("	%s, %d\n", p.Name, p.Stats.Frags)
		}
	}*/
	
	fmt.Printf("Teams count : %d\n", len(teams))
	
	event_path := w.DataPath + "/events.json"
	
	events := InitEvents(event_path)
	
	w.Data.Events = events
	
	fmt.Printf("Events count : %d\n", len(events))
	
	m_count := 0
	
	for _, evt := range events {
		m_count = m_count + len(evt.Matches)
	}
	
	fmt.Printf("Matches count : %d\n", m_count)
	
	if len(w.Data.MissingTeams) > 0 {
		//fetching missings teams
		fmt.Println("Missing teams!")
		for _, tId := range w.Data.MissingTeams {
			
			team := &csgoscrapper.Team{Name:"NotSet", TeamId: tId}
			team.LoadTeam()
			
			fmt.Printf("Fecthing team [%d] \n", team.TeamId)
			
			w.Data.Teams = append(w.Data.Teams, team)
			
		}
		
		w.Data.MissingTeams = []int{}
		
	}
	
	w.Running = false
}

func (w *WatcherState) StartBot()  {
	
	d := time.Minute * 5
	
	for {
		fmt.Printf("Sleeping %f minutes...\n", d.Minutes())
		time.Sleep(d)
		
		w.Running = true
		
		fmt.Println("Updating events and matches from HLTV.org")
		
		page := csgoscrapper.GetEventsPage()
		
		pc, _ := page.LoadPage()
		
		w.Data.Events = pc.UpdateEvents(w.Data.Events)
		
		//need to save events
		fmt.Println("Saving events...")
		event_path := w.DataPath + "/events.json"
		
		csgoscrapper.SaveEvents(w.Data.Events, event_path)
		
		//saving teams
		
		fmt.Println("Saving teams...")
		
		team_path := w.DataPath + "/teams.json"
		csgoscrapper.SaveTeams(w.Data.Teams, team_path)
		
		w.Running = false
	}
	
}