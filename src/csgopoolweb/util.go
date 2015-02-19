package csgopoolweb
/*
import (
	"csgoscrapper"
	"fmt"
)
//can be remove

func (w *WebServerState) GetEventById(id int) *csgoscrapper.Event {
	
	event := csgoscrapper.GetEventById(w.Data.Events, id)
	
	return event
	
}

func (w *WebServerState) GetTeamById(id int) *csgoscrapper.Team {
	
	team := csgoscrapper.GetTeamById(w.Data.Teams, id)
	
	if team == nil {
		
		//w.AddMissingTeam(id)
		state.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", id))
	
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
}*/
