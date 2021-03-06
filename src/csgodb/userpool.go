package csgodb

import (
	"database/sql"
	"time"
)

type UserPool struct {
	PoolId int
	DivisionId int
	UserId int
	PlayerId int
	CreatedOn time.Time
}



type MetaPool struct {
	UserPool
	Username string
	Division Division
	Player Player
	Points int
}

func InsertPoolChoices(db *sql.DB, choices []*UserPool) {
	
	query := "INSERT INTO users_pools (division_id, user_id, player_id, created_on) VALUES (?, ?, ?, ?)"
	
	now := time.Now()
	
	for _, choice := range choices {
		db.Exec(query, choice.DivisionId, choice.UserId, choice.PlayerId, now)
	}
	
}


func GetMetaPoolsByUserWithoutPoints(db *sql.DB, userId int) []*MetaPool {
	pools := []*MetaPool{}
	query := `SELECT up.pool_id, up.division_id, up.user_id, up.player_id, u.username, p.player_name, d.division_name 
		FROM users_pools up 
		JOIN users u ON u.user_id = up.user_id 
		JOIN players p ON p.player_id = up.player_id 
		JOIN divisions d ON d.division_id = up.division_id 
		WHERE up.user_id = ? `
	rows, _ := db.Query(query, userId)
	
	for rows.Next() {
		pool := &MetaPool{}
		rows.Scan(&pool.PoolId, &pool.DivisionId, &pool.UserId, &pool.PlayerId, &pool.Username, &pool.Player.Name, &pool.Division.Name)
		pool.Player.PlayerId = pool.PlayerId
		pool.Division.DivisionId = pool.DivisionId
		pools = append(pools, pool)
	}
	
	return pools
}

func GetMetaPoolsByUser(db *sql.DB, userId int) []*MetaPool {
	
	pools := []*MetaPool{}
	
	query := `SELECT up.pool_id, up.division_id, up.user_id, up.player_id, u.username,  p.player_name, d.division_name, SUM(pp.points) as points
	FROM users_pools up 
	JOIN users u ON u.user_id = up.user_id 
	JOIN players p ON p.player_id = up.player_id 
	JOIN divisions d ON d.division_id = up.division_id 
	LEFT JOIN players_points pp ON pp.player_id = up.player_id 
	LEFT JOIN matches m ON m.match_id = pp.match_id
	WHERE up.user_id = ? AND (DATE(up.created_on) <= m.match_date)
	GROUP BY up.pool_id 
	ORDER BY division_id ASC`

	rows, _ := db.Query(query, userId)
	
	for rows.Next() {
		pool := &MetaPool{}
		rows.Scan(&pool.PoolId, &pool.DivisionId, &pool.UserId, &pool.PlayerId, &pool.Username, &pool.Player.Name, &pool.Division.Name, &pool.Points)
		pool.Player.PlayerId = pool.PlayerId
		pool.Division.DivisionId = pool.DivisionId
		pools = append(pools, pool)
	}
	
	return pools
	
}

func GetMetaPools(db *sql.DB) []*MetaPool {
	
	pools := []*MetaPool{}
	
	query := "SELECT up.pool_id, up.division_id, up.user_id, up.player_id, u.username,  p.player_name, d.division_name "
	query += "FROM users_pools up "
	query += "JOIN users u ON u.user_id = up.user_id "
	query += "JOIN players p ON p.player_id = up.player_id "
	query += "JOIN divisions d ON d.division_id = up.division_id "
	rows, _ := db.Query(query)
	
	for rows.Next() {
		pool := &MetaPool{}
		rows.Scan(&pool.PoolId, &pool.DivisionId, &pool.UserId, &pool.PlayerId, &pool.Username, &pool.Player.Name, &pool.Division.Name)
		pool.Player.PlayerId = pool.PlayerId
		pool.Division.DivisionId = pool.DivisionId
		pools = append(pools, pool)
	}
	
	return pools
	
}

func GetAllUserPools(db *sql.DB) []*UserPool {
	pools := []*UserPool{}
	
	query := "SELECT pool_id, division_id, user_id, player_id FROM users_pools"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		pool := &UserPool{}
		rows.Scan(&pool.PoolId, &pool.DivisionId, &pool.UserId, &pool.PlayerId)
		pools = append(pools, pool)
	}
	
	return pools
}

func GetPoolsByUser(db *sql.DB, userId int) []*UserPool {
	pools := []*UserPool{}
	
	query := "SELECT pool_id, division_id, user_id, player_id FROM users_pools WHERE user_id = ?"
	rows, _ := db.Query(query, userId)
	
	for rows.Next() {
		pool := &UserPool{}
		rows.Scan(&pool.PoolId, &pool.DivisionId, &pool.UserId, &pool.PlayerId)
		pools = append(pools, pool)
	}
	
	return pools
}
