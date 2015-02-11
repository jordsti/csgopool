package csgopoolweb

import (
	"csgoscrapper"
	"fmt"
)

func (w *WebServerState) GetEventById(id int) *csgoscrapper.Event {
	
	event := csgoscrapper.GetEventById(w.Data.Events, id)
	
	return event
	
}

func (w *WebServerState) GetTeamById(id int) *csgoscrapper.Team {
	
	team := csgoscrapper.GetTeamById(w.Data.Teams, id)
	
	if team == nil {
		
		//w.AddMissingTeam(id)
		fmt.Printf("Team [%d] not found, fetching this team!\n", id)
	
		newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: id}
		newTeam.LoadTeam()
		
		w.Data.Teams = append(w.Data.Teams, newTeam)
	
		return newTeam
		
	} else {
		return team
	}
	
}

func (w *WebServerState) GetMatchById(id int) *csgoscrapper.Match {
	match := csgoscrapper.GetMatchById(w.Data.Events, id)
	return match
}

func (w *WebServerState) AddMissingTeam(teamId int) {
	
	for _, id := range w.Data.MissingTeams {
		if id == teamId {
			return
		}
	}
	
	w.Data.MissingTeams = append(w.Data.MissingTeams, teamId)
	
}