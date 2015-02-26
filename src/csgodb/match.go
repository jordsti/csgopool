package csgodb

import (
	"hltvscrapper"
	"eseascrapper"
	"database/sql"
	"time"
	"fmt"
)

type MatchTeam struct {
	TeamId int
	Name string
	Score int
}

type Match struct {
	MatchId int
	EventId int
	Team1 MatchTeam
	Team2 MatchTeam
	Map string
	Date time.Time
	PlayerStats []PlayerStat
	Source int
	SourceId int
	PoolStatus int
	SourceName string
}

func GetLastMatch(db *sql.DB) *Match {
	
	match := &Match{MatchId: 0}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, source, source_id, pool_status FROM matches ORDER BY match_id DESC LIMIT 1"
	
	rows, _ := db.Query(query)
	
	for rows.Next() {
		rows.Scan(&match.MatchId, &match.Team1.TeamId, &match.Team1.Score, &match.Team2.TeamId, &match.Team2.Score, &match.Map, &match.EventId, &match.Date, &match.Source, &match.SourceId, &match.PoolStatus)
	}
	
	return match
}


func GetMatchById(db *sql.DB, matchId int) *Match {
	
	match := &Match{MatchId: 0}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, source, source_id, pool_status FROM matches WHERE match_id = ?"
	
	rows, _ := db.Query(query, matchId)
	
	for rows.Next() {
		rows.Scan(&match.MatchId, &match.Team1.TeamId, &match.Team1.Score, &match.Team2.TeamId, &match.Team2.Score, &match.Map, &match.EventId, &match.Date, &match.Source, &match.SourceId, &match.PoolStatus)
	}
	
	return match
}

func UpdateMatchPoolStatus(db *sql.DB, matchId int, poolStatus int) {
	query := "UPDATE matches SET pool_status = ? WHERE match_id = ?"
	db.Exec(query, poolStatus, matchId)
}

func IsMatchExists(db *sql.DB, matchId int) bool {
	
	match := Match{MatchId: 0}
	query := "SELECT match_id FROM matches WHERE match_id = ?"
	rows, _ := db.Query(query, matchId)
	
	for rows.Next() {
		rows.Scan(&match.MatchId)
	}
	
	if match.MatchId != 0 {
		return true
	}
	
	return false
	
}

func IsMatchIn(matches []*Match, matchId int) bool {
	
	for _, m := range matches {
		if m.MatchId == matchId {
			return true
		}
	}
	
	return false
	
}

func IsSourceMatchIn(matches []*Match, source int, sourceId int) bool {
	for _, m := range matches {
		if m.Source == source && m.SourceId == sourceId {
			return true
		}
	}
	
	return false
}

func GetAllMatches(db *sql.DB) []*Match {
	
	matches := []*Match{}
	
	query := `SELECT m.match_id, m.team1_id, t1.team_name, m.team1_score, m.team2_id, t2.team_name, m.team2_score, m.map, m.event_id, m.match_date, m.source, m.source_id, m.pool_status, s.source_name 
			FROM matches m
			JOIN teams t1 ON t1.team_id = m.team1_id 
			JOIN teams t2 ON t2.team_id = m.team2_id
			JOIN sources s ON s.source_id = m.source
			ORDER BY match_date DESC`
	rows, _ := db.Query(query)
	
	for rows.Next() {
		m := &Match{}
		
		rows.Scan(&m.MatchId, &m.Team1.TeamId, &m.Team1.Name, &m.Team1.Score, &m.Team2.TeamId, &m.Team2.Name, &m.Team2.Score, &m.Map, &m.EventId, &m.Date, &m.Source, &m.SourceId, &m.PoolStatus, &m.SourceName)
		matches = append(matches, m)
	}
	
	return matches
}


func GetMatchesByEventId(db *sql.DB, eventId int) []*Match {
	
	matches := []*Match{}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, source, source_id, pool_status FROM matches WHERE event_id = ? ORDER BY match_id DESC"
	rows, _ := db.Query(query, eventId)
	
	for rows.Next() {
		m := &Match{}
		
		rows.Scan(&m.MatchId, &m.Team1.TeamId, &m.Team1.Score, &m.Team2.TeamId, &m.Team2.Score, &m.Map, &m.EventId, &m.Date, &m.Source, &m.SourceId, &m.PoolStatus)
		matches = append(matches, m)
	}
	
	return matches
}

func GetMatchesByDate(db *sql.DB, date time.Time) []*Match {
	
	matches := []*Match{}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, source, source_id, pool_status FROM matches WHERE match_date >= ? ORDER BY match_id DESC"
	rows, _ := db.Query(query, date)
	
	for rows.Next() {
		m := &Match{}
		
		rows.Scan(&m.MatchId, &m.Team1.TeamId, &m.Team1.Score, &m.Team2.TeamId, &m.Team2.Score, &m.Map, &m.EventId, &m.Date, &m.Source, &m.SourceId, &m.PoolStatus)
		matches = append(matches, m)
	}
	
	return matches
}

func GetMatchBySource(db *sql.DB, source int, sourceId int) *Match {
	m := &Match{}
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, source, source_id, pool_status FROM matches WHERE source = ? AND source_id = ?"
	rows, _ := db.Query(query, source, sourceId)
	
	for rows.Next() {
		rows.Scan(&m.MatchId, &m.Team1.TeamId, &m.Team1.Score, &m.Team2.TeamId, &m.Team2.Score, &m.Map, &m.EventId, &m.Date, &m.Source, &m.SourceId, &m.PoolStatus)
	}
	
	return m
}

func ImportHltvMatch(db *sql.DB, m *hltvscrapper.Match) *Match {
	query := "INSERT INTO matches (source, source_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	
	date := time.Date(m.Date.Year, time.Month(m.Date.Month), m.Date.Day, 0, 0, 0, 0, time.Local)
	db.Exec(query, HltvSource, m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score, m.Map, m.Event.EventId, date, 0)
	
	_m := GetMatchBySource(db, HltvSource, m.MatchId)
	_m.ImportHltvStats(db, m.PlayerStats)
	
	return _m
}

func ImportEseaMatch(db *sql.DB, m *eseascrapper.Match) *Match {
	
	query := "INSERT INTO matches (source, source_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	
	date := time.Date(m.Date.Year, time.Month(m.Date.Month), m.Date.Day, 0, 0, 0, 0, time.Local)
	db.Exec(query, EseaSource, m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score, m.Map, 0, date, 0)
	
	_m := GetMatchBySource(db, EseaSource, m.MatchId)
	
	_m.ImportEseaStats(db, m.PlayerStats)
	
	return _m
}

func (m *Match) ImportHltvStats(db *sql.DB, stats []*hltvscrapper.MatchPlayerStat) {
	query := "INSERT INTO matches_stats (match_id, team_id, player_id, frags, assists, deaths, kdratio, kddelta) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	//team id must be fixed with real id not source
	//same for player id
	for _, ms := range stats {
		db.Exec(query, m.MatchId, ms.TeamId, ms.PlayerId, ms.Frags, ms.Assists, ms.Deaths, ms.KDRatio, ms.KDDelta)
	}
}

func (m *Match) ImportEseaStats(db *sql.DB, stats []*eseascrapper.PlayerMatchStat) {
	query := "INSERT INTO matches_stats (match_id, team_id, player_id, frags, assists, deaths, kdratio, kddelta) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	
	for _, ms := range stats {
			db.Exec(query, m.MatchId, ms.TeamId, ms.PlayerId, ms.Frags, ms.Assists, ms.Deaths, ms.KDRatio, ms.KDDelta)
	}
}

func ImportMatch(db *sql.DB, m hltvscrapper.Match) {
	

	query := "INSERT INTO matches (match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	
	date := time.Date(m.Date.Year, time.Month(m.Date.Month), m.Date.Day, 0, 0, 0, 0, time.Local)
	_, _ = db.Exec(query, m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score, m.Map, m.Event.EventId, date, 0)
	
	ImportMatchesStats(db, m.MatchId, m.PlayerStats)
	
}


func ImportMatches(db *sql.DB, matches []*hltvscrapper.Match) {
	
	for _, m := range matches {
		query := "INSERT INTO matches (match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		
		date := time.Date(m.Date.Year, time.Month(m.Date.Month), m.Date.Day, 0, 0, 0, 0, time.Local)
		_, _ = db.Exec(query, m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score, m.Map, m.Event.EventId, date, 0)
		
		ImportMatchesStats(db, m.MatchId, m.PlayerStats)
	}
}

func ImportMatchesStats(db *sql.DB, matchId int, stats []*hltvscrapper.MatchPlayerStat) {
	query := "INSERT INTO matches_stats (match_id, team_id, player_id, frags, headshots, assists, deaths, kdratio, kddelta) VALUES (?, ?, ?, ?, ?, ?,  ?, ?, ?)"
	for _, s := range stats {
		
		if stats != nil {
			_, err := db.Exec(query, matchId, s.TeamId, s.PlayerId, s.Frags, s.Headshots, s.Assists, s.Deaths, s.KDRatio, s.KDDelta)
			if err != nil {
				fmt.Printf("SQL Error: %v\n", err)
			}
		} else {
			fmt.Printf("Match stats nil [%d]\n", matchId)
		}
	}
}
