package csgodb

import (
	"database/sql"
	"csgoscrapper"
)

type Player struct {
	PlayerId int
	Name string
}

func ImportPlayers(db *sql.DB, players []csgoscrapper.Player) {
	
	query := "INSERT INTO players (player_id, player_name) VALUES (?, ?)"
	
	for _, pl := range players {
		db.Exec(query, pl.PlayerId, pl.Name)
	}
	
}