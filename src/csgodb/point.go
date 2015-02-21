package csgodb

import (
	"database/sql"
)


func AddPoint(db *sql.DB, matchId int, playerId int, points int) {
	query := "INSERT INTO players_points (match_id, player_id, points) VALUES(?, ?, ?) "
	db.Exec(query, matchId, playerId, points)
}

/*
little query
SELECT * FROM players_points pt
JOIN matches m ON m.match_id = pt.match_id 
JOIN players p ON p.player_id = pt.player_id
JOIN matches_stats ms ON ms.match_id = pt.match_id AND ms.player_id = pt.player_id
JOIN teams t ON t.team_id = ms.team_id
JOIN `events` e ON e.event_id = m.event_id
WHERE m.match_id = 19161
ORDER BY m.match_id DESC

*/