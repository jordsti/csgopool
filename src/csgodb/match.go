package csgodb

import (
	"csgoscrapper"
	"database/sql"
	"time"
	"fmt"
)

type MatchTeam struct {
	TeamId int
	Score int
}

type Match struct {
	MatchId int
	EventId int
	Team1 MatchTeam
	Team2 MatchTeam
	Map string
	Date time.Time
}


func GetAllMatches(db *sql.DB) []*Match {
	
	matches := []*Match{}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date FROM matches"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		m := &Match{}
		
		rows.Scan(&m.MatchId, &m.Team1.TeamId, &m.Team1.Score, &m.Team2.TeamId, &m.Team2.Score, &m.Map, &m.EventId, &m.Date)
		matches = append(matches, m)
	}
	
	return matches
}


func GetMatchesByEventId(db *sql.DB, eventId int) []*Match {
	
	matches := []*Match{}
	
	query := "SELECT match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date FROM matches WHERE event_id = ?"
	rows, _ := db.Query(query, eventId)
	
	for rows.Next() {
		m := &Match{}
		
		rows.Scan(&m.MatchId, &m.Team1.TeamId, &m.Team1.Score, &m.Team2.TeamId, &m.Team2.Score, &m.Map, &m.EventId, &m.Date)
		matches = append(matches, m)
	}
	
	return matches
}

func ImportMatches(db *sql.DB, matches []csgoscrapper.Match) {
	
	for _, m := range matches {
		query := "INSERT INTO matches (match_id, team1_id, team1_score, team2_id, team2_score, map, event_id, match_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		_, _ = db.Exec(query, m.MatchId, m.Team1.TeamId, m.Team1.Score, m.Team2.TeamId, m.Team2.Score, m.Map, m.EventId, 0)
		
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
