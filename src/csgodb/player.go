package csgodb

import (
	"database/sql"
	"csgoscrapper"
	"eseascrapper"
	"time"
	"fmt"
)

type GeneralStat struct {
	MatchesPlayed int
	Frags int
	//Headshots int
	Assists int
	Deaths int
	AvgKDRatio float32
	AvgKDDelta float32
	AvgFrags float32
}

type Player struct {
	PlayerId int
	Name string
	EseaId int
	HltvId int
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
	KDRatio float32
	Points int
}

func (pl *Player) P() *PlayerWithStat {
	plP := &PlayerWithStat{}
	plP.PlayerId = pl.PlayerId
	plP.Name = pl.Name
	
	return plP
}

func GetAllPlayersWithStat(db *sql.DB) []*PlayerWithStat {
	players := []*PlayerWithStat{}
	
	query := "SELECT p.player_id, p.player_name, SUM(ms.frags), SUM(ms.deaths), AVG(ms.kdratio), COUNT(ms.match_stat_id) "
	query += "FROM players p "
	query += "JOIN matches_stats ms ON ms.player_id = p.player_id "
	query += "GROUP BY player_id ORDER BY p.player_name"
	
	rows, _ := db.Query(query)
	for rows.Next() {
		player := &PlayerWithStat{}
		rows.Scan(&player.PlayerId, &player.Name, &player.Stat.Frags, &player.Stat.Deaths, &player.Stat.AvgKDRatio, &player.Stat.MatchesPlayed)
		players = append(players, player)
	}
	
	return players
}

func GetPlayerMatchStats(db *sql.DB, playerId int) []*PlayerMatchStat {
	stats := []*PlayerMatchStat{}
	
	query := "SELECT m.match_id, m.match_date, m.team1_id, t1.team_name, m.team1_score, "
	query += "m.team2_id, t2.team_name, m.team2_score, ms.frags, ms.kdratio, pp.points "
	query += "FROM matches m "
	query += "JOIN teams t1 ON t1.team_id = m.team1_id "
	query += "JOIN teams t2 ON t2.team_id = m.team2_id "
	query += "JOIN matches_stats ms ON ms.match_id = m.match_id "
	query += "LEFT JOIN players_points pp ON pp.player_id = ms.player_id AND pp.match_id = ms.match_id "
	query += "WHERE ms.player_id = ? ORDER BY m.match_date DESC "
	
	rows, _ := db.Query(query, playerId)
	
	for rows.Next() {
		stat := &PlayerMatchStat{}
		
		rows.Scan(&stat.MatchId, &stat.Date, &stat.Team1.TeamId, &stat.Team1.Name, &stat.TeamScore1, &stat.Team2.TeamId, &stat.Team2.Name, &stat.TeamScore2, &stat.Frags, &stat.KDRatio, &stat.Points)
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
	query += " SUM(ms.assists), SUM(ms.deaths), AVG(ms.kdratio), AVG(ms.kddelta), AVG(ms.frags)"
	query += "FROM players p "
	query += "JOIN players_teams pt ON pt.player_id = p.player_id "
	query += "JOIN teams t ON t.team_id = pt.team_id "
	query += "JOIN matches_stats ms ON ms.player_id = p.player_id AND ms.team_id = pt.team_id "
	query += "WHERE p.player_id = ? GROUP BY player_id"
	
	rows, _ := db.Query(query, playerId)
	for rows.Next() {
		rows.Scan(&pl.PlayerId, &pl.Name, &pl.Stat.MatchesPlayed, &pl.Stat.Frags, &pl.Stat.Assists, &pl.Stat.Deaths, &pl.Stat.AvgKDRatio, &pl.Stat.AvgKDDelta, &pl.Stat.AvgFrags)
	}
	
	return pl
}

func GetPlayerById(db *sql.DB, playerId int) *Player {
	
	player := &Player{PlayerId: 0}
	
	query := "SELECT player_id, player_name, esea_id, hltv_id FROM players WHERE player_id = ?"
	rows, _ := db.Query(query, playerId)
	
	for rows.Next() {
		rows.Scan(&player.PlayerId, &player.Name, &player.EseaId, &player.HltvId)
	}
	
	if player.PlayerId != 0 {
		return player
	}
	
	return nil
}

func GetAllPlayers(db *sql.DB) []*Player {
	
	players := []*Player{}
	
	query := "SELECT player_id, player_name, esea_id, hltv_id FROM players ORDER BY player_name"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		
		player := &Player{}
		rows.Scan(&player.PlayerId, &player.Name, &player.EseaId, &player.HltvId)
		players = append(players, player)
	}
	
	return players
	
}

func GetPlayersWithStatByTeamId(db *sql.DB, teamId int) []*PlayerWithStat {
	players := []*PlayerWithStat{}
	
	query := "SELECT p.player_id, p.player_name, COUNT(ms.match_stat_id), SUM(ms.frags), "
	query += " SUM(ms.assists), SUM(ms.deaths), AVG(ms.kdratio), AVG(ms.kddelta), AVG(ms.frags)"
	query += "FROM players p "
	query += "JOIN players_teams pt ON pt.player_id = p.player_id "
	query += "JOIN matches_stats ms ON ms.player_id = p.player_id AND ms.team_id = pt.team_id "
	query += "WHERE pt.team_id = ? GROUP BY player_id"
	
	rows, _ := db.Query(query, teamId)
	
	for rows.Next() {
		pl := &PlayerWithStat{}
		rows.Scan(&pl.PlayerId, &pl.Name, &pl.Stat.MatchesPlayed, &pl.Stat.Frags, &pl.Stat.Assists, &pl.Stat.Deaths, &pl.Stat.AvgKDRatio, &pl.Stat.AvgKDDelta, &pl.Stat.AvgFrags)
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

func IsSourcePlayerIn(players []*Player, source int, sourceId int) bool {
	
	for _, pl := range players {
		
		if source == EseaSource && pl.EseaId == sourceId {
			return true
		} else if source == HltvSource && pl.HltvId == sourceId {
			return true
		}
	}
	
	return false
}

func FindPlayerByName(players []*Player, name string) *Player {
	for _, pl := range players {
		if pl.Name == name {
			return pl
		}
	}
	
	return nil
}

func GetPlayerByName(db *sql.DB, name string) *Player {
	player := &Player{}
	query := "SELECT player_id, player_name, esea_id, hltv_id FROM players WHERE player_name = ?"
	rows, _ := db.Query(query, name)
	
	for rows.Next() {
		rows.Scan(&player.PlayerId, &player.Name, &player.EseaId, &player.HltvId)
	}
	
	return player
}

func (p *Player) UpdateSourceId(db *sql.DB) {
	query := "UPDATE players SET hltv_id = ?, esea_id = ? WHERE player_id = ?"
	db.Exec(query, p.HltvId, p.EseaId, p.PlayerId)
}

//need to be sure that this player doesnt exists
func ImportHltvPlayer(db *sql.DB, player csgoscrapper.Player) *Player {
	query := "INSERT INTO players (player_name, hltv_id) VALUES (?, ?)"
	
	db.Exec(query, player.Name, player.PlayerId)
	
	pl := GetPlayerByName(db, player.Name)
	return pl
}

func ImportEseaPlayer(db *sql.DB, player eseascrapper.Player) *Player {
	query := "INSERT INTO players (player_name, esea_id) VALUES (?, ?)"
	
	_, err := db.Exec(query, player.Name, player.PlayerId)
	
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	
	pl := GetPlayerByName(db, player.Name)
	return pl
}

//this will go deprecated
func ImportPlayer(db *sql.DB, player csgoscrapper.Player) {
	
	query := "INSERT INTO players (player_id, player_name) VALUES (?, ?)"
	
	db.Exec(query, player.PlayerId, player.Name)
	
}
//same as up here
func ImportPlayers(db *sql.DB, players []csgoscrapper.Player) {
	
	query := "INSERT INTO players (player_id, player_name) VALUES (?, ?)"
	
	for _, pl := range players {
		db.Exec(query, pl.PlayerId, pl.Name)
	}
	
}