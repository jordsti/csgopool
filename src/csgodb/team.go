package csgodb

import (
	"csgoscrapper"
	"database/sql"
)


type Team struct {
	TeamId int
	Name string
	
	Players []Player
}


func AddPlayer(db *sql.DB, teamId int, playerId int) {
	query := "INSERT INTO players_teams (team_id, player_id) VALUES (?, ?)"
	
	db.Exec(query, teamId, playerId)
}

func IsTeamExists(teams []*Team, teamId int) bool {
	for _, t := range teams {
		if t.TeamId == teamId { return true }
	}
	
	return false
}

func ImportTeams(db *sql.DB, teams []*csgoscrapper.Team) {
	//this is for initial import only !!
	
	for _, team := range teams {
		
		stmt, _ := db.Prepare("INSERT INTO teams (team_id, team_name) VALUES (?, ?)")
		stmt.Exec(team.TeamId, team.Name)
		defer stmt.Close()
	}
	
}

func GetAllTeams(db *sql.DB) []*Team {
	
	teams := []*Team{}
	
	rows, _ := db.Query("SELECT team_id, team_name FROM teams")

	for rows.Next() {
		team := &Team{}
		rows.Scan(&team.TeamId, &team.Name)
		teams = append(teams, team)
	}
	
	return teams
}

func GetTeamById(db *sql.DB, teamId int) Team {
	
	team := Team{TeamId: 0, Name:""}
	
	rows, _ := db.Query("SELECT team_id, team_name FROM teams WHERE team_id = ?")
		
	for rows.Next() {
		rows.Scan(&team.TeamId, &team.Name)
	}
	
	return team
}