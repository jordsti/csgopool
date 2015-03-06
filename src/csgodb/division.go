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
	query := "DELETE FROM users_pools"
	db.Exec(query)
	
	query = "DELETE FROM divisions_players"
	db.Exec(query)
	
	query = "DELETE FROM divisions"
	db.Exec(query)
}

func DivisionById(divisions []*Division, divisionId int) *Division {
	
	for _, div := range divisions {
		if div.DivisionId == divisionId {
			return div
		}
	}
	
	return nil
}

func DivisionCount(db *sql.DB) int {
	
	query := "SELECT division_id FROM divisions"
	rows, _ := db.Query(query)
	it := 0
	for rows.Next() {
		it++
	}
	
	return it
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

func GetAllDivisionsWithPlayer(db *sql.DB) []*Division {
	
	divs := []*Division{}
	
	query := "SELECT dp.division_id, dp.player_id, d.division_name, p.player_name FROM divisions_players dp "
	query += "JOIN players p ON p.player_id = dp.player_id "
	query += "JOIN divisions d ON d.division_id = dp.division_id "
	query += "ORDER BY d.division_id "
	
	rows, _ := db.Query(query)
	
	currentDiv := &Division{DivisionId: 0}
	
	divs = append(divs, currentDiv)
	it := 0
	for rows.Next() {
		d_id := 0
		d_name := ""
		pl := &Player{}
		rows.Scan(&d_id, &pl.PlayerId, &d_name, &pl.Name)
		
		if d_id != currentDiv.DivisionId {
			if currentDiv.DivisionId == 0 {
				//first division
				currentDiv.DivisionId = d_id
				currentDiv.Name = d_name
				
			} else {
				//pushing new division throw the slice
				//and create a new one
				currentDiv = &Division{DivisionId: d_id, Name: d_name}
				divs = append(divs, currentDiv)
			}
		}
		
		currentDiv.Players = append(currentDiv.Players, pl)
		it++
	}
	
	if it == 0 {
		divs = []*Division{}
	}
	
	return divs
}

func (d *Division) IsPlayerIn(playerId int) bool {
	
	for _, pl := range d.Players {
		if pl.PlayerId == playerId {
			return true
		}
	}
	
	return false
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