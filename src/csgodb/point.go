package csgodb

import (
	"database/sql"
	"time"
)

type MatchPointStat struct {
	MatchId int
	MatchDate time.Time
	PlayerId int
	PlayerName string
	TeamId int
	TeamName string
	Frags int
	Headshots int
	KDRatio float32
	Points int
}

type PlayerPoints struct {
	PlayerId int
	Name string
	Matches int
	Frags int
	Headshots int
	KDRatio float32
	KDDelta float32
	Points int
}

type UserPoints struct {
	UserId int
	Name string
	Points int
}

type PlayerDivisionPoints struct {
	PlayerId int
	Name string
	Points int
}

type DivisionPoints struct {
	Players []*PlayerDivisionPoints
	DivisionId int
	Name string
}

func GetDivisionsPoints(db *sql.DB) []*DivisionPoints {
	points := []*DivisionPoints{}
	
	query := `SELECT dp.division_id, d.division_name, p.player_id, p.player_name, SUM(pp.points) as points
			FROM divisions_players dp 
			JOIN divisions d ON dp.division_id = d.division_id  
			JOIN players p ON p.player_id = dp.player_id
			LEFT JOIN players_points pp ON pp.player_id = dp.player_id
			GROUP BY dp.player_id
			ORDER BY d.division_id `
	
	rows, _ := db.Query(query)
	currentDiv := &DivisionPoints{}
	
	for rows.Next() {
		d_id := 0
		d_name := ""
		pl := &PlayerDivisionPoints{}
		rows.Scan(&d_id, &d_name, &pl.PlayerId, &pl.Name, &pl.Points)
		
		if currentDiv.DivisionId != d_id {
			currentDiv = &DivisionPoints{DivisionId: d_id, Name: d_name}
			points = append(points, currentDiv)
		}
		
		currentDiv.Players = append(currentDiv.Players, pl)
		
	}
	
	return points
}

func GetPlayersPoint(db *sql.DB) []*PlayerPoints {
	
	points := []*PlayerPoints{}
	query := `SELECT p.player_id, p.player_name, COUNT(ms.match_stat_id), SUM(ms.frags), SUM(ms.headshots), AVG(ms.kdratio), AVG(ms.kddelta), SUM(pt.points) as points FROM players_points pt
				JOIN players p ON p.player_id = pt.player_id
				JOIN matches_stats ms ON ms.match_id = pt.match_id AND ms.player_id = pt.player_id
				GROUP BY player_id
				ORDER BY points DESC`
	
	rows, _ := db.Query(query)
	
	for rows.Next() {
		point := &PlayerPoints{}
		rows.Scan(&point.PlayerId, &point.Name, &point.Matches, &point.Frags, &point.Headshots, &point.KDRatio, &point.KDDelta, &point.Points)
		points = append(points, point)
	}
	
	return points
}

func GetUserPoints(db *sql.DB) []*UserPoints {
	points := []*UserPoints{}
	
	query := `SELECT u.user_id, u.username, SUM(pt.points) as points FROM users u
			LEFT JOIN users_pools up ON up.user_id = u.user_id 
			LEFT JOIN players_points pt ON up.player_id = pt.player_id
			JOIN matches m ON m.match_id = pt.match_id
			WHERE (DATE(up.created_on) <= m.match_date)
			GROUP BY up.user_id ORDER BY points DESC`
	rows, _ := db.Query(query)
	
	for rows.Next() {
		point := &UserPoints{}
		rows.Scan(&point.UserId, &point.Name, &point.Points)
		points = append(points, point)
	}
	
	return points
}

func GetMatchPoints(db *sql.DB, matchId int) []*MatchPointStat {
	stats := []*MatchPointStat{}
	
	query := `SELECT m.match_date, p.player_id, p.player_name, t.team_id, t.team_name, ms.frags, ms.kdratio, pt.points FROM players_points pt 
	JOIN matches m ON m.match_id = pt.match_id 
	JOIN players p ON p.player_id = pt.player_id 
	JOIN matches_stats ms ON ms.match_id = pt.match_id AND ms.player_id = pt.player_id 
	JOIN teams t ON t.team_id = ms.team_id 
	JOIN events e ON e.event_id = m.event_id 
	WHERE m.match_id = ? 
	ORDER BY pt.points DESC`
	
	rows, _ := db.Query(query, matchId)

	
	for rows.Next() {
		stat := &MatchPointStat{}
		rows.Scan(&stat.MatchDate, &stat.PlayerId, &stat.PlayerName, &stat.TeamId, &stat.TeamName, &stat.Frags, &stat.KDRatio, &stat.Points)
		
		stat.MatchId = matchId
		stats = append(stats, stat)
	}
	
	return stats
}

func AddPoint(db *sql.DB, matchId int, playerId int, points int) {
	query := "INSERT INTO players_points (match_id, player_id, points) VALUES(?, ?, ?) "
	db.Exec(query, matchId, playerId, points)
}