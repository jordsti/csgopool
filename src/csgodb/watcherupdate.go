package csgodb

import (
	"database/sql"
	"time"
)

type WatcherUpdate struct {
	UpdateId int
	Time time.Time
}

func InsertWatcherUpdate(db *sql.DB) {
	
	date := time.Now()
	
	query := "INSERT INTO watcher_update (update_time) VALUES (?)"
	db.Exec(query, date)
}

func GetLastUpdate(db *sql.DB) *WatcherUpdate {
	update := &WatcherUpdate{}
	
	query := "SELECT update_id, update_time FROM watcher_update ORDER BY update_time DESC LIMIT 1"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		rows.Scan(&update.UpdateId, &update.Time)
	}
	
	return update
}

