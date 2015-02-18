package csgopool

import (
	"fmt"
	"csgoscrapper"
	"csgodb"
	"time"
	"os"
	"database/sql"
)


var watcher *WatcherState

type GameData struct {
	Events []*csgoscrapper.Event
	Teams []*csgoscrapper.Team
}

type WatcherState struct {
	
	DataPath string
	Data GameData //to remove mysql is used instead todo
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

func (w *WatcherState) InitialImport(db *sql.DB) {
	
	w.Log.Info("Initial import from HLTV.org")
	
	teams := InitTeams()
	events := InitEvents()
	
	w.Log.Info("Scanning matches for missing team")
	for _, evt := range events {
		for _, m := range evt.Matches {
			team1 := csgoscrapper.GetTeamById(teams, m.Team1.TeamId)
			team2 := csgoscrapper.GetTeamById(teams, m.Team2.TeamId)
			
			if team1 == nil {
				w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team1.TeamId))
				newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team1.TeamId}
				newTeam.LoadTeam()
				teams = append(teams, newTeam)	
			}
			
			if team2 == nil {
				w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team2.TeamId))
				newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team2.TeamId}
				newTeam.LoadTeam()
				teams = append(teams, newTeam)	
			}
		}
	}
	
	w.Log.Info("Inserting data into database")
	
	csgodb.ImportTeams(db, teams)
	
	for _, t := range teams {
		
		csgodb.ImportPlayers(db, t.Players)
		for _, p := range t.Players {
			csgodb.AddPlayerToTeam(db, t.TeamId, p.PlayerId)
		}
		
	}
	
	csgodb.ImportEvents(db, events)
	
	for _, evt := range events {
		if len(evt.Matches) > 0 {
			
			for _, m := range evt.Matches {
				csgodb.ImportMatch(db, m)
			}
			
		}
	}
	
	w.Log.Info("hltv import terminated")
	w.Log.Info(fmt.Sprintf("%d teams, %d events imported !", len(teams), len(events)))
	
}

func (w *WatcherState) LoadData() {
	
	w.Running = true
	
	//init database
	//todo
	//config file
	
	db_path := w.DataPath + "/db.json"
	csgodb.Db.LoadConfig(db_path)
	
	if len(csgodb.Db.Username) == 0 {
		w.Log.Error("Invalid Database config, please edit db.json!")
		csgodb.Db.SaveConfig(db_path)
		os.Exit(-2)
	}
	
	db, err := csgodb.Db.Open()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	csgodb.InitTables(db)
	
	teams := csgodb.GetAllTeams(db)
	
	if len(teams) == 0 {
		w.Log.Info("No team found !")
		w.Log.Info("Starting initial import...")
		w.InitialImport(db)
	}
	

	if csgodb.UsersCount(db) == 0 {
		w.Log.Info("No user found, Creating default PoolMaster Account")
		
		passwd := RandomString(12)
		
		csgodb.CreateUser(db, "poolmaster", passwd, "poolmaster@localhost", PoolMaster)
		
		w.Log.Info("PoolMaster Account Created !")
		w.Log.Info(fmt.Sprintf("poolmaster:%s", passwd))
	}
	
	
	db.Close()
	w.Running = false
}

func (w *WatcherState) StartBot()  {
	//rework this!
	d := time.Minute * 5

	w.Log.Info("Starting watcher Bot")
	for {
		w.Log.Info(fmt.Sprintf("Sleeping %f minutes...", d.Minutes()))
		time.Sleep(d)
		
		w.Running = true
		
		//updating last events
		
		db, _ := csgodb.Db.Open()
		
		lastEvent := csgodb.GetLastEvent(db)
		players := csgodb.GetAllPlayers(db)
		events := csgodb.GetAllEvents(db)
		//check the last 5 events if possible

		for key, evt := range events {
			if key <  5 {
					w.Log.Info(fmt.Sprintf("Update event [%d]-%s", evt.EventId, evt.Name))
					matches := csgodb.GetMatchesByEventId(db, evt.EventId)
					teams := []*csgoscrapper.Team{}
					
					url := csgoscrapper.GetEventMatches(evt.EventId)
					//todo
					//error handling
					pc, _ := url.LoadPage()
					
					event_matches := pc.ParseMatches()
					new_matches := []*csgoscrapper.Match{}
					
					for _, m := range event_matches {
						
						if !csgodb.IsMatchIn(matches, m.MatchId) {
							w.Log.Info(fmt.Sprintf("Match [%d] not in event [%d], retrieving player stats", m.MatchId, evt.EventId))
							m.GetMatchStats()
							
							new_matches = append(new_matches, m)
						}
						
					}
					
					//missing teams
					for _, m := range new_matches {
					team1 := csgodb.GetTeamById(db, m.Team1.TeamId)
					team2 := csgodb.GetTeamById(db, m.Team2.TeamId)
					
						if team1.TeamId == 0 {
							w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team1.TeamId))
							newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team1.TeamId}
							newTeam.LoadTeam()
							teams = append(teams, newTeam)	
						}
						
						if team2.TeamId == 0 {
							w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team2.TeamId))
							newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team2.TeamId}
							newTeam.LoadTeam()
							teams = append(teams, newTeam)	
						}
					
					}
					
					//importing teams
					csgodb.ImportTeams(db, teams)
					for _, t := range teams {
						
						for _, pl := range t.Players {
							if !csgodb.IsPlayerIn(players, pl.PlayerId) {
								csgodb.ImportPlayer(db, pl)
							}
						}
		
						for _, p := range t.Players {
							csgodb.AddPlayerToTeam(db, t.TeamId, p.PlayerId)
						}
						
					}
					//importing matches
					csgodb.ImportMatches(db, new_matches)
					
					if len(new_matches) > 0 {
						evt.Tick(db)
					}
			} else { break }
		}
		//checking for new events
		
		w.Log.Debug(fmt.Sprintf("events : %d", len(events)))
		url := csgoscrapper.GetEventsPage()
		teams := []*csgoscrapper.Team{}
		
		pc, _ := url.LoadPage()
		
		n_events := pc.ParseEventsWithoutMatches()
		new_events := []*csgoscrapper.Event{}
		//reloading
		players = csgodb.GetAllPlayers(db)

		for _, evt := range n_events {
			
			if !csgodb.IsEventIn(events, evt.EventId) && evt.EventId > lastEvent.EventId {
				w.Log.Info(fmt.Sprintf("Event [%d] not in database", evt.EventId))
				evt.LoadAllMatches()
				
				if len(evt.Matches) > 0 {
					new_events = append(new_events, evt)
					
					for _, m := range evt.Matches {
						team1 := csgodb.GetTeamById(db, m.Team1.TeamId)
						team2 := csgodb.GetTeamById(db, m.Team2.TeamId)
				
						if team1.TeamId == 0 {
							w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team1.TeamId))
							newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team1.TeamId}
							newTeam.LoadTeam()
							teams = append(teams, newTeam)	
						}
						
						if team2.TeamId == 0 {
							w.Log.Info(fmt.Sprintf("Team [%d] not found, fetching this team", m.Team2.TeamId))
							newTeam := &csgoscrapper.Team{Name: "NotSet", TeamId: m.Team2.TeamId}
							newTeam.LoadTeam()
							teams = append(teams, newTeam)	
						}
					}
					
					n_matches := []*csgoscrapper.Match{}
					
					for _, m := range evt.Matches {
						n_matches = append(n_matches, &m)
					}
					
					csgodb.ImportTeams(db, teams)
					
					for _, t := range teams {
						for _, pl := range t.Players {
							if csgodb.IsPlayerIn(players, pl.PlayerId) {
								csgodb.ImportPlayer(db, pl)
							}
						}
	
						for _, p := range t.Players {
							csgodb.AddPlayerToTeam(db, t.TeamId, p.PlayerId)
						}
					
					}
					//importing matches
					csgodb.ImportMatches(db, n_matches)
				} else {
					w.Log.Info("0 matches found, probably a too old event..")
					break
				}
			}
			
			csgodb.ImportEvents(db, new_events)
		}
		
		db.Close()
		
		w.Running = false
	}
	
}