package csgopool

import (
	"fmt"
	"csgoscrapper"
	"time"
)


var watcher *WatcherState

type GameData struct {
	Events []*csgoscrapper.Event
	Teams []*csgoscrapper.Team
}

type WatcherState struct {
	
	DataPath string
	Data GameData
	Running bool
	Users Users
	Log *csgoscrapper.LoggerState
}



func NewWatcher(dataPath string) *WatcherState {
	
	state := &WatcherState{DataPath: dataPath}
	state.Running = false
	state.Log = &csgoscrapper.LoggerState{LogPath: dataPath+"/watcher.log", Level:3}
	watcher = state
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
	
	w.Log.Info(fmt.Sprintf("Teams count : %d", len(teams)))
	
	event_path := w.DataPath + "/events.json"
	
	events := InitEvents(event_path)
	
	w.Data.Events = events
	
	w.Log.Info(fmt.Sprintf("Events count : %d", len(events)))
	
	m_count := 0
	w.Log.Info("Scanning matches for not existing team")
	for _, evt := range events {
		m_count = m_count + len(evt.Matches)
		for _, m := range evt.Matches {
			team1 := csgoscrapper.GetTeamById(w.Data.Teams, m.Team1.TeamId)
			team2 := csgoscrapper.GetTeamById(w.Data.Teams, m.Team2.TeamId)
			
			if team1 == nil {
				w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team1.TeamId))
				newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team1.TeamId}
				newTeam.LoadTeam()
				w.Data.Teams = append(w.Data.Teams, newTeam)	
			}
			
			if team2 == nil {
				w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team2.TeamId))
				newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team2.TeamId}
				newTeam.LoadTeam()
				w.Data.Teams = append(w.Data.Teams, newTeam)	
			}
		}
	}
	
	csgoscrapper.SaveTeams(w.Data.Teams, team_path)
	
	w.Log.Info(fmt.Sprintf("Matches count : %d", m_count))
	
	
	users := &w.Users
	users.LoadUsers(w.DataPath + "users.json")
	
	if len(w.Users.Users) == 0 {
		w.Log.Info("No user found, Creating default PoolMaster Account")
		
		//todo
		//random password show the credentials in the console output
		
		passwd := RandomString(12)
		
		w.Users.CreateUser("poolmaster", passwd, "poolmaster@localhost", PoolMaster)
		
		w.Log.Info("PoolMaster Account Created !")
		w.Log.Info(fmt.Sprintf("poolmaster:%s", passwd))
	}
	
	w.Running = false
}

func (w *WatcherState) StartBot()  {
	
	d := time.Minute * 5
	
	for {
		w.Log.Info(fmt.Sprintf("Sleeping %f minutes...", d.Minutes()))
		time.Sleep(d)
		
		w.Running = true
		
		w.Log.Info("Updating events and matches from HLTV.org")
		
		page := csgoscrapper.GetEventsPage()
		
		pc, _ := page.LoadPage()
		
		w.Data.Events = pc.UpdateEvents(w.Data.Events)
		
		//need to save events
		w.Log.Info("Saving events...")
		event_path := w.DataPath + "/events.json"
		
		csgoscrapper.SaveEvents(w.Data.Events, event_path)
		
		//checking for missing team
		
		for _, evt := range w.Data.Events {
			for _, m := range evt.Matches {
				team1 := csgoscrapper.GetTeamById(w.Data.Teams, m.Team1.TeamId)
				team2 := csgoscrapper.GetTeamById(w.Data.Teams, m.Team2.TeamId)
				
				if team1 == nil {
					w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team1.TeamId))
					newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team1.TeamId}
					newTeam.LoadTeam()
					w.Data.Teams = append(w.Data.Teams, newTeam)	
				}
				
				if team2 == nil {
					w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team2.TeamId))
					newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team2.TeamId}
					newTeam.LoadTeam()
					w.Data.Teams = append(w.Data.Teams, newTeam)
				}
			}
		}
		
		//saving teams
		
		w.Log.Info("Saving teams...")
		
		team_path := w.DataPath + "/teams.json"
		csgoscrapper.SaveTeams(w.Data.Teams, team_path)
		
		w.Running = false
	}
	
}