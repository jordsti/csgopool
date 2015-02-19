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
	Headshots int
	Assists int
	Deaths int
	KDRatio float32
	KDDelta int
	PlayerName string
}

func (m *Match) FetchStats(db *sql.DB) {
	
	query := "SELECT ms.match_stat_id, ms.match_id, ms.team_id, ms.player_id, ms.frags, ms.headshots, ms.assists, ms.deaths, ms.kdratio, ms.kddelta, p.player_name FROM matches_stats ms JOIN players p ON p.player_id = ms.player_id WHERE match_id = ? ORDER BY ms.kddelta DESC"
	rows, _ := db.Query(query, m.MatchId)
	
	for rows.Next() {
		stat := PlayerStat{}
		rows.Scan(&stat.MatchStatId, &stat.MatchId, &stat.TeamId, &stat.PlayerId, &stat.Frags, &stat.Headshots, &stat.Assists, &stat.Deaths, &stat.KDRatio, &stat.KDDelta, &stat.PlayerName)
		m.PlayerStats = append(m.PlayerStats, stat)
	}
}
