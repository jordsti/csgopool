package csgodb

import (
	"csgoscrapper"
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
	PoolStatus int
}

func GetLastMatch(db *sql.DB) *Match {
	
	match := &Match{MatchId: 0}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status FROM matches ORDER BY match_id DESC LIMIT 1"
	
	rows, _ := db.Query(query)
	
	for rows.Next() {
		rows.Scan(&match.MatchId, &match.Team1.TeamId, &match.Team1.Score, &match.Team2.TeamId, &match.Team2.Score, &match.Map, &match.EventId, &match.Date, &match.PoolStatus)
	}
	
	return match
}


func GetMatchById(db *sql.DB, matchId int) *Match {
	
	match := &Match{MatchId: 0}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status FROM matches WHERE match_id = ?"
	
	rows, _ := db.Query(query, matchId)
	
	for rows.Next() {
		rows.Scan(&match.MatchId, &match.Team1.TeamId, &match.Team1.Score, &match.Team2.TeamId, &match.Team2.Score, &match.Map, &match.EventId, &match.Date, &match.PoolStatus)
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

func GetAllMatches(db *sql.DB) []*Match {
	
	matches := []*Match{}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status FROM matches"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		m := &Match{}
		
		rows.Scan(&m.MatchId, &m.Team1.TeamId, &m.Team1.Score, &m.Team2.TeamId, &m.Team2.Score, &m.Map, &m.EventId, &m.Date, &m.PoolStatus)
		matches = append(matches, m)
	}
	
	return matches
}


func GetMatchesByEventId(db *sql.DB, eventId int) []*Match {
	
	matches := []*Match{}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status FROM matches WHERE event_id = ? ORDER BY match_id DESC"
	rows, _ := db.Query(query, eventId)
	
	for rows.Next() {
		m := &Match{}
		
		rows.Scan(&m.MatchId, &m.Team1.TeamId, &m.Team1.Score, &m.Team2.TeamId, &m.Team2.Score, &m.Map, &m.EventId, &m.Date, &m.PoolStatus)
		matches = append(matches, m)
	}
	
	return matches
}

func ImportMatch(db *sql.DB, m csgoscrapper.Match) {
	

	query := "INSERT INTO matches (match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	
	date := time.Date(m.Date.Year, time.Month(m.Date.Month), m.Date.Day, 0, 0, 0, 0, time.Local)
	_, _ = db.Exec(query, m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score, m.Map, m.EventId, date, 0)
	
	ImportMatchesStats(db, m.MatchId, m.PlayerStats)
	
}


func ImportMatches(db *sql.DB, matches []*csgoscrapper.Match) {
	
	for _, m := range matches {
		query := "INSERT INTO matches (match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date, pool_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		
		date := time.Date(m.Date.Year, time.Month(m.Date.Month), m.Date.Day, 0, 0, 0, 0, time.Local)
		_, _ = db.Exec(query, m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score, m.Map, m.EventId, date, 0)
		
		ImportMatchesStats(db, m.MatchId, m.PlayerStats)
	}
}

func ImportMatchesStats(db *sql.DB, matchId int, stats []*csgoscrapper.MatchPlayerStat) {
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
