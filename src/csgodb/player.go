package csgodb

import (
	"database/sql"
	"csgoscrapper"
)

type Player struct {
	PlayerId int
	Name string
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

func GetPlayersByTeamId(db *sql.DB, teamId int) []*Player {
	players := []*Player{}
	
	query := "SELECT p.player_id, p.player_name FROM players p JOIN players_teams pt ON pt.player_id = p.player_id WHERE pt.team_id = ?"
	
	rows, _ := db.Query(query, teamId)
	
	for rows.Next() {
		pl := &Player{}
		rows.Scan(&pl.PlayerId, &pl.Name)
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