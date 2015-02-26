package csgodb

import (
	"hltvscrapper"
	"eseascrapper"
	"database/sql"
	"fmt"
)


type Team struct {
	TeamId int
	Name string
	EseaId int
	HltvId int
	
	Players []*Player
}

type TeamP struct {
	Team
	Players []*PlayerWithStat
	PlayersCount int
	MatchesCount int
}

func GetTeamMatches(db *sql.DB, teamId int) []*Match {
	matches := []*Match{}
	
	query := "SELECT m.match_id, m.team1_id, t1.team_name, m.team1_score, m.team2_id, t2.team_name, m.team2_score, m.map, m.match_date "
	query += "FROM matches m "
	query += "JOIN teams t1 ON t1.team_id = m.team1_id "
	query += "JOIN teams t2 ON t2.team_id = m.team2_id "
	query += "WHERE m.team1_id = ? OR m.team2_id = ? ORDER BY m.match_date DESC"
	
	rows, err := db.Query(query, teamId, teamId)
	
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	
	for rows.Next() {
		match := &Match{}
		rows.Scan(&match.MatchId, &match.Team1.TeamId, &match.Team1.Name, &match.Team1.Score, &match.Team2.TeamId, &match.Team2.Name, &match.Team2.Score, &match.Map, &match.Date)
		matches = append(matches, match)
	}
	
	return matches
}

func GetTeamsWithCount(db *sql.DB) []*TeamP {
	teams := []*TeamP{}
	query := "SELECT t.team_id, t.team_name, COUNT(pt.player_id), (SELECT COUNT(m.match_id) FROM matches m WHERE m.team1_id = t.team_id OR m.team2_id = t.team_id) as played_matches FROM teams t JOIN players_teams pt ON pt.team_id = t.team_id GROUP BY team_id ORDER BY played_matches DESC"
	
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

func ImportTeams(db *sql.DB, teams []*hltvscrapper.Team) {
	//this is for initial import only !!
	
	/*for _, team := range teams {
		
		//query := "INSERT INTO teams (team_id, team_name, esea_id, hltv_id) VALUES (?, ?, ?, ?)"
		//todo
	}*/
	
}

func FindTeamByName(teams []*Team, name string) *Team {
	for _, t := range teams {
		if t.Name == name {
			return t
		}
	}
	
	return nil
}

func ImportHltvTeam(db *sql.DB, team *hltvscrapper.Team) *Team {
	query := "INSERT INTO teams (team_name, hltv_id) VALUES (?, ?)"
	db.Exec(query, team.Name, team.TeamId)
	return GetTeamByHltvId(db, team.TeamId)
}

func ImportEseaTeam(db *sql.DB, team *eseascrapper.Team) *Team {
	query := "INSERT INTO teams (team_name, esea_id) VALUES (?, ?)"
	db.Exec(query, team.Name, team.TeamId)
	return GetTeamByEseaId(db, team.TeamId)
}

func GetTeamByEseaId(db *sql.DB, eseaId int) *Team {
	team := &Team{}
	
	query := "SELECT team_id, team_name, esea_id, hltv_id FROM teams WHERE esea_id = ?"
	rows, _ := db.Query(query, eseaId)
	
	for rows.Next() {
		rows.Scan(&team.TeamId, &team.Name, &team.EseaId, &team.HltvId)
	}
	
	return team
}

func GetTeamByHltvId(db *sql.DB, hltvId int) *Team {
	team := &Team{}
	
	query := "SELECT team_id, team_name, esea_id, hltv_id FROM teams WHERE hltv_id = ?"
	rows, _ := db.Query(query, hltvId)
	
	for rows.Next() {
		rows.Scan(&team.TeamId, &team.Name, &team.EseaId, &team.HltvId)
	}
	
	return team
}

func GetTeamByName(db *sql.DB, name string) *Team {
	team := &Team{}
	
	query := "SELECT team_id, team_name, esea_id, hltv_id FROM teams WHERE team_name = ?"
	rows, _ := db.Query(query, name)
	
	for rows.Next() {
		rows.Scan(&team.TeamId, &team.Name, &team.EseaId, &team.HltvId)
	}
	
	return team
}

func (t *Team) UpdateSourceId(db *sql.DB) {
	query := "UPDATE teams SET esea_id = ?, hltv_id = ? WHERE team_id = ?"
	db.Exec(query, t.EseaId, t.HltvId, t.TeamId)
}


func GetAllTeams(db *sql.DB) []*Team {
	
	teams := []*Team{}
	
	rows, _ := db.Query("SELECT team_id, team_name, esea_id, hltv_id FROM teams")

	for rows.Next() {
		team := &Team{}
		rows.Scan(&team.TeamId, &team.Name, &team.EseaId, &team.HltvId)
		teams = append(teams, team)
	}
	
	return teams
}

func GetTeamById(db *sql.DB, teamId int) Team {
	
	team := Team{TeamId: 0, Name:""}
	
	rows, err := db.Query("SELECT team_id, team_name, esea_id, hltv_id FROM teams WHERE team_id = ?", teamId)
	
	if err != nil {
		fmt.Printf("%v\n", err)
	}
		
	for rows.Next() {
		rows.Scan(&team.TeamId, &team.Name, &team.EseaId, &team.HltvId)
	}
	
	return team
}