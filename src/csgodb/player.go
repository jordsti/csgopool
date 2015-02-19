package csgodb

import (
	"database/sql"
	"csgoscrapper"
	"time"
)

type GeneralStat struct {
	MatchesPlayed int
	Frags int
	Headshots int
	Assists int
	Deaths int
	AvgKDRatio float32
	AvgKDDelta float32
	AvgFrags float32
	AvgHeadshots float32
}

type Player struct {
	PlayerId int
	Name string
}

type PlayerWithStat struct {
	Player
	Stat GeneralStat
}

type PlayerTeamStat struct {
	TeamId int
	Name string
	MatchesCount int
}

type PlayerMatchStat struct {
	MatchId int
	Date time.Time
	Team1 Team
	TeamScore1 int
	Team2 Team
	TeamScore2 int
	Map string
	Frags int
	Headshots int
	KDRatio float32
}

func (pl *Player) P() *PlayerWithStat {
	plP := &PlayerWithStat{}
	plP.PlayerId = pl.PlayerId
	plP.Name = pl.Name
	
	return plP
}

func GetAllPlayersWithStat(db *sql.DB) []*PlayerWithStat {
	players := []*PlayerWithStat{}
	
	query := "SELECT p.player_id, p.player_name, SUM(ms.frags), SUM(ms.headshots), SUM(ms.deaths), AVG(ms.kdratio), COUNT(ms.match_stat_id) "
	query += "FROM players p "
	query += "JOIN matches_stats ms ON ms.player_id = p.player_id "
	query += "GROUP BY player_id ORDER BY p.player_name"
	
	rows, _ := db.Query(query)
	for rows.Next() {
		player := &PlayerWithStat{}
		rows.Scan(&player.PlayerId, &player.Name, &player.Stat.Frags, &player.Stat.Headshots, &player.Stat.Deaths, &player.Stat.AvgKDRatio, &player.Stat.MatchesPlayed)
		players = append(players, player)
	}
	
	return players
}

func GetPlayerMatchStats(db *sql.DB, playerId int) []*PlayerMatchStat {
	stats := []*PlayerMatchStat{}
	
	query := "SELECT m.match_id, m.match_date, m.team1_id, t1.team_name, m.team1_score, "
	query += "m.team2_id, t2.team_name, m.team2_score, ms.frags, ms.headshots, ms.kdratio "
	query += "FROM matches m "
	query += "JOIN teams t1 ON t1.team_id = m.team1_id "
	query += "JOIN teams t2 ON t2.team_id = m.team2_id "
	query += "JOIN matches_stats ms ON ms.match_id = m.match_id "
	query += "WHERE ms.player_id = ? ORDER BY m.match_date DESC"
	
	rows, _ := db.Query(query, playerId)
	
	for rows.Next() {
		stat := &PlayerMatchStat{}
		
		rows.Scan(&stat.MatchId, &stat.Date, &stat.Team1.TeamId, &stat.Team1.Name, &stat.TeamScore1, &stat.Team2.TeamId, &stat.Team2.Name, &stat.TeamScore2, &stat.Frags, &stat.Headshots, &stat.KDRatio)
		stats = append(stats, stat)
	}
	
	return stats
}

func GetPlayerTeamStats(db *sql.DB, playerId int) []*PlayerTeamStat {
	
	teams := []*PlayerTeamStat{}
	
	query := "SELECT t.team_id, t.team_name, COUNT(ms.match_stat_id) FROM matches_stats ms "
	query += "JOIN teams t ON t.team_id = ms.team_id "
	query += "WHERE ms.player_id = ? GROUP BY team_id "
	
	rows, _ := db.Query(query, playerId)
	
	for rows.Next() {
		team := &PlayerTeamStat{}
		rows.Scan(&team.TeamId, &team.Name, &team.MatchesCount)
		teams = append(teams, team)
	}
	
	return teams
}

func GetPlayerWithStatById(db *sql.DB, playerId int) *PlayerWithStat {
	pl := &PlayerWithStat{}
	query := "SELECT p.player_id, p.player_name, COUNT(ms.match_stat_id), SUM(ms.frags), "
	query += "SUM(ms.headshots), SUM(ms.assists), SUM(ms.deaths), AVG(ms.kdratio), AVG(ms.kddelta), AVG(ms.frags), AVG(ms.headshots)"
	query += "FROM players p "
	query += "JOIN players_teams pt ON pt.player_id = p.player_id "
	query += "JOIN teams t ON t.team_id = pt.team_id "
	query += "JOIN matches_stats ms ON ms.player_id = p.player_id AND ms.team_id = pt.team_id "
	query += "WHERE p.player_id = ? GROUP BY player_id"
	
	rows, _ := db.Query(query, playerId)
	for rows.Next() {
		rows.Scan(&pl.PlayerId, &pl.Name, &pl.Stat.MatchesPlayed, &pl.Stat.Frags, &pl.Stat.Headshots, &pl.Stat.Assists, &pl.Stat.Deaths, &pl.Stat.AvgKDRatio, &pl.Stat.AvgKDDelta, &pl.Stat.AvgFrags, &pl.Stat.AvgHeadshots)
	}
	
	return pl
}

func GetPlayerById(db *sql.DB, playerId int) *Player {
	
	player := &Player{PlayerId: 0}
	
	query := "SELECT player_id, player_name FROM players WHERE player_id = ?"
	rows, _ := db.Query(query, playerId)
	
	for rows.Next() {
		rows.Scan(&player.PlayerId, &player.Name)
	}
	
	if player.PlayerId != 0 {
		return player
	}
	
	return nil
}

func GetAllPlayers(db *sql.DB) []*Player {
	
	players := []*Player{}
	
	query := "SELECT player_id, player_name FROM players ORDER BY player_name"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		
		player := &Player{}
		rows.Scan(&player.PlayerId, &player.Name)
		players = append(players, player)
	}
	
	return players
	
}

func GetPlayersWithStatByTeamId(db *sql.DB, teamId int) []*PlayerWithStat {
	players := []*PlayerWithStat{}
	
	query := "SELECT p.player_id, p.player_name, COUNT(ms.match_stat_id), SUM(ms.frags), "
	query += "SUM(ms.headshots), SUM(ms.assists), SUM(ms.deaths), AVG(ms.kdratio), AVG(ms.kddelta), AVG(ms.frags), AVG(ms.headshots)"
	query += "FROM players p "
	query += "JOIN players_teams pt ON pt.player_id = p.player_id "
	query += "JOIN matches_stats ms ON ms.player_id = p.player_id AND ms.team_id = pt.team_id "
	query += "WHERE pt.team_id = ? GROUP BY player_id"
	
	rows, _ := db.Query(query, teamId)
	
	for rows.Next() {
		pl := &PlayerWithStat{}
		rows.Scan(&pl.PlayerId, &pl.Name, &pl.Stat.MatchesPlayed, &pl.Stat.Frags, &pl.Stat.Headshots, &pl.Stat.Assists, &pl.Stat.Deaths, &pl.Stat.AvgKDRatio, &pl.Stat.AvgKDDelta, &pl.Stat.AvgFrags, &pl.Stat.AvgHeadshots)
		players = append(players, pl)
	}
	
	return players
}

func IsPlayerIn(players []*Player, playerId int) bool {
	
	for _, pl := range players {
		if pl.PlayerId == playerId {
			return true
		}
	}
	
	return false
}

func ImportPlayer(db *sql.DB, player csgoscrapper.Player) {
	
	query := "INSERT INTO players (player_id, player_name) VALUES (?, ?)"
	
	db.Exec(query, player.PlayerId, player.Name)
	
}

func ImportPlayers(db *sql.DB, players []csgoscrapper.Player) {
	
	query := "INSERT INTO players (player_id, player_name) VALUES (?, ?)"
	
	for _, pl := range players {
		db.Exec(query, pl.PlayerId, pl.Name)
	}
	
}