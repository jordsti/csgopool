package csgodb

import (
	"database/sql"
	"csgoscrapper"
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

func GetAllPlayers(db *sql.DB) []*Player {
	
	players := []*Player{}
	
	query := "SELECT player_id, player_name FROM players"
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