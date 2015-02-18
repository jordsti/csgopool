package csgodb

import (
	"csgoscrapper"
	"database/sql"
	"fmt"
)


type Team struct {
	TeamId int
	Name string
	
	Players []*Player
}

type TeamP struct {
	Team
	Players []*PlayerWithStat
	PlayersCount int
	MatchesCount int
}

func GetTeamsWithCount(db *sql.DB) []*TeamP {
	teams := []*TeamP{}
	query := "SELECT t.team_id, t.team_name, COUNT(pt.player_id), (SELECT COUNT(m.match_id) FROM matches m WHERE m.team1_id = t.team_id OR m.team2_id = t.team_id) as played_matches FROM teams t JOIN players_teams pt ON pt.team_id = t.team_id GROUP BY t.team_id"
	
	rows, _ := db.Query(query)
	
	for rows.Next() {
		tp := &TeamP{}
		rows.Scan(&tp.TeamId, &tp.Name, &tp.PlayersCount, &tp.MatchesCount)
		teams = append(teams, tp)
	}
	
	return teams
}

func (t *Team) P() *TeamP {
	team := &TeamP{}
	
	team.TeamId = t.TeamId
	team.Name = t.Name
	
	return team
}

func (t *TeamP) FetchPlayers(db *sql.DB) {
	
	t.Players = GetPlayersWithStatByTeamId(db, t.TeamId)
}

func AddPlayerToTeam(db *sql.DB, teamId int, playerId int) {
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
	
	rows, err := db.Query("SELECT team_id, team_name FROM teams WHERE team_id = ?", teamId)
	
	if err != nil {
		fmt.Printf("%v\n", err)
	}
		
	for rows.Next() {
		rows.Scan(&team.TeamId, &team.Name)
	}
	
	return team
}