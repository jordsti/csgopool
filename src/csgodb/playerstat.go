package csgodb

import (
	"database/sql"
)

type PlayerStat struct {
	MatchStatId int
	MatchId int
	TeamId int
	PlayerId int
	Frags int
	//Headshots int
	Assists int
	Deaths int
	KDRatio float32
	KDDelta int
	PlayerName string
	Points int
}

func (m *Match) FetchStats(db *sql.DB) {
	
	query := `SELECT ms.match_stat_id, ms.match_id, ms.team_id, ms.player_id, ms.frags, ms.assists, ms.deaths, ms.kdratio, ms.kddelta, p.player_name, pp.points 
	FROM matches_stats ms 
	JOIN players p ON p.player_id = ms.player_id 
	LEFT JOIN players_points pp ON pp.player_id = ms.player_id AND pp.match_id = ms.match_id
	WHERE ms.match_id = ? 
	ORDER BY ms.kddelta DESC`
	rows, _ := db.Query(query, m.MatchId)
	
	for rows.Next() {
		stat := PlayerStat{}
		rows.Scan(&stat.MatchStatId, &stat.MatchId, &stat.TeamId, &stat.PlayerId, &stat.Frags, &stat.Assists, &stat.Deaths, &stat.KDRatio, &stat.KDDelta, &stat.PlayerName, &stat.Points)
		m.PlayerStats = append(m.PlayerStats, stat)
	}
}
