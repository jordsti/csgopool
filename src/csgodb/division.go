package csgodb

import (
	"database/sql"
)

type Division struct {
	DivisionId int
	Name string
	Players []*Player
}

func ClearPool(db *sql.DB) {
	query := "DELETE FROM divisions_players"
	db.Exec(query)
	
	query = "DELETE FROM divisions"
	db.Exec(query)
}

func GetLastDivisionId(db *sql.DB) int {
	
	div_id := 0
	query := "SELECT division_id FROM divisions ORDER BY division_id DESC LIMIT 1"
	
	rows, _ := db.Query(query)
	
	for rows.Next() {
		rows.Scan(&div_id)
	}
	
	return div_id
}

func AddDivision(db *sql.DB, name string) *Division {
	
	query := "INSERT INTO divisions (division_name) VALUES (?)"
	
	db.Exec(query, name)
	
	div := &Division{}
	
	div.DivisionId = GetLastDivisionId(db)
	div.Name = name
	
	return div
}

func GetAllDivisions(db *sql.DB) []*Division {
	divs := []*Division{}
	
	query := "SELECT division_id, division_name FROM divisions"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		div := &Division{}
		rows.Scan(&div.DivisionId, &div.Name)
		divs = append(divs, div)
	}
	
	return divs
}

func (d *Division) AddPlayer(db *sql.DB, playerId int) {
	query := "INSERT INTO divisions_players (player_id, division_id) VALUES (?, ?)"
	db.Exec(query, playerId, d.DivisionId)
}

func (d *Division) FetchPlayers(db *sql.DB) {
	players  := []*Player{}
	query := "SELECT p.player_id, p.player_name FROM players p "
	query += "JOIN divisions_players dp ON dp.player_id = p.player_id "
	query += "WHERE dp.division_id = ?"
	
	rows, _ := db.Query(query, d.DivisionId)
	for rows.Next() {
		player := &Player{}
		rows.Scan(&player.PlayerId, &player.Name)
		players = append(players, player)
	}
	
	d.Players = players
}