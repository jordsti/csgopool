package csgodb

import (
	"database/sql"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"fmt"
)

type PlayerTeam struct {
	TeamId int
	PlayerId int
}

type MatchStat struct {
	MatchId int
	TeamId int
	PlayerId int
	Frags int
	Assists int
	Deaths int
	KDRatio float32
	KDDelta int
	Source int
	SourceId int
}

type Snapshot struct {
	Matches []*Match
	Players []*Player
	Teams []*Team
	MatchesStats []*MatchStat
	Events []*Event
	PlayerTeam []*PlayerTeam
}

func GenerateSnapshot(db *sql.DB) *Snapshot {
	snapshot := &Snapshot{}
	
	//events
	query := "SELECT event_id, event_name FROM events"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		event := &Event{}
		rows.Scan(&event.EventId, &event.Name)
		snapshot.Events = append(snapshot.Events, event)
	}
	
	//players
	query = "SELECT player_id, player_name, esea_id, hltv_id FROM players"
	rows, _ = db.Query(query)
	for rows.Next() {
		player := &Player{}
		rows.Scan(&player.PlayerId, &player.Name, &player.EseaId, &player.HltvId)
		snapshot.Players = append(snapshot.Players, player)
	}
	
	//teams
	query = "SELECT team_id, team_name, esea_id, hltv_id FROM teams"
	rows, _ = db.Query(query)
	for rows.Next() {
		team := &Team{}
		rows.Scan(&team.TeamId, &team.Name, &team.EseaId, &team.HltvId)
		snapshot.Teams = append(snapshot.Teams, team)
	}
	
	//player -> team
	query = "SELECT player_id, team_id FROM players_teams"
	rows, _ = db.Query(query)
	for rows.Next() {
		pteam := &PlayerTeam{}
		rows.Scan(&pteam.PlayerId, &pteam.TeamId)
		snapshot.PlayerTeam = append(snapshot.PlayerTeam, pteam)
	}
	
	//matches
	query = "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, source, souce_id FROM matches"
	rows, _ = db.Query(query)
	for rows.Next() {
		match := &Match{}
		rows.Scan(&match.MatchId, &match.Team1.TeamId, &match.Team1.Score, &match.Team2.TeamId, &match.Team2.Score, &match.Map, &match.EventId, &match.Date, &match.Source, &match.SourceId)
		snapshot.Matches = append(snapshot.Matches, match)
	}
	
	//matches stats
	query = "SELECT match_id, team_id, player_id, frags, assists, deaths, kdratio, kddelta FROM matches_stats"
	rows, _ = db.Query(query)
	for rows.Next() {
		stat := &MatchStat{}
		rows.Scan(&stat.MatchId, &stat.TeamId, &stat.PlayerId, &stat.Frags, &stat.Assists, &stat.Deaths, &stat.KDRatio, &stat.KDDelta)
		snapshot.MatchesStats = append(snapshot.MatchesStats, stat)
	}
	
	return snapshot
}

func (s *Snapshot) Save(path string) {
	b, err := json.MarshalIndent(s, "", "	")
	if err != nil {
		fmt.Println("Error while converting Snapshot to JSON")
	}
	
	err = ioutil.WriteFile(path, b, 0644)
	
	if err != nil {
		fmt.Println("Error while saving Snapshot")
	}
}

func (s *Snapshot) Parse(jsonData []byte) {
	
	err := json.Unmarshal(jsonData, s)
	if err != nil {
		fmt.Println("Error while parsing Snapshot")
	}
}

func (s *Snapshot) ImportFromURL(db *sql.DB, url string) {

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error while trying to fetch a Snapshot [1]")
	}
	
	body, err2 := ioutil.ReadAll(resp.Body)
	
	if err2 != nil {
		fmt.Println("Error while trying to fetch a Snapshot [2]")
	}
	
	if err == nil && err2 == nil {
		s.Parse(body)
		s.Import(db)
	}
}

func (s *Snapshot) Import(db *sql.DB) {
	//players
	for _, pl := range s.Players {
		query := "INSERT INTO players (player_id, player_name, esea_id, hltv_id) VALUES (?, ?, ?, ?)"
		db.Exec(query, pl.PlayerId, pl.Name, pl.EseaId, pl.HltvId)
	}
	
	//teams
	for _, t := range s.Teams {
		query := "INSERT INTO teams (team_id, team_name, esea_id, hltv_id) VALUES (?, ?, ?, ?)"
		db.Exec(query, t.TeamId, t.Name, t.EseaId, t.HltvId)
	}
	
	//players_teams
	for _, pt := range s.PlayerTeam {
		query := "INSERT INTO players_teams (team_id, player_id) VALUES (?, ?)"
		db.Exec(query, pt.TeamId, pt.PlayerId)
	}
	
	//events
	for _, evt := range s.Events {
		query := "INSERT INTO events (event_id, event_name) VALUES (?, ?)"
		db.Exec(query, evt.EventId, evt.Name)
	}
	
	//matches
	for _, m := range s.Matches {
		query := "INSERT INTO matches (match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, source, source_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		db.Exec(query, m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score, m.Map, m.EventId, m.Date, m.Source, m.SourceId)
	}
	
	//match stats
	for _, ms := range s.MatchesStats {
		query := "INSERT INTO matches_stats (match_id, team_id, player_id, frags, assists, deaths, kdratio, kddelta) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		db.Exec(query, ms.MatchId, ms.TeamId, ms.PlayerId, ms.Frags, ms.Assists, ms.Deaths, ms.KDRatio, ms.KDDelta)
	}
}