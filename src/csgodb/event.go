package csgodb
import (
	"database/sql"
	"csgoscrapper"
	"fmt"
)

type Event struct {
	EventId int
	Name string
	Matches []*Match
}

func ImportEvents(db *sql.DB, events []*csgoscrapper.Event) {
	
	for _, evt := range events {
		
		stmt, err := db.Prepare("INSERT INTO events (event_id, event_name) VALUES (?, ?)")
		
		_, err = stmt.Query(evt.EventId, evt.Name)
		
		defer stmt.Close()
		
		if err != nil {
			fmt.Printf("SQL Error: %v\n", err)
		}
		
	}
	
}