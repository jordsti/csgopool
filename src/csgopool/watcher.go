package csgopool

import (
	"fmt"
	"csgoscrapper"
)


type GameData struct {
	Events []*csgoscrapper.Event
	Teams []*csgoscrapper.Team
}

type WatcherState struct {
	
	DataPath string
	Data GameData
}



func NewWatcher(dataPath string) *WatcherState {
	
	state := &WatcherState{DataPath: dataPath}
	return state
}

func (w *WatcherState) LoadData() {
	
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
	
}