package csgodb
import (
	"database/sql"
	"csgoscrapper"
	"fmt"
	"time"
)

type Event struct {
	EventId int
	Name string
	Matches []*Match
	MatchesCount int
	LastChange time.Time
}

func IsEventIn(events []*Event, eventId int) bool {
	for _, e := range events {
		if e.EventId == eventId {
			return true
		}
	}
	
	return false
}

func TickEvent(db *sql.DB, eventId int) {
	last_change := time.Now()
	
	query := "UPDATE events SET last_change = ? WHERE event_id = ?"
	db.Exec(query, last_change, eventId)
}

func (e *Event) Tick(db *sql.DB) {
	
	last_change := time.Now()
	
	query := "UPDATE events SET last_change = ? WHERE event_id = ?"
	db.Exec(query, last_change, e.EventId)
	
}

func IsEventExists(db *sql.DB, eventId int) bool {
	
	event := Event{EventId: 0}
	query := "SELECT event_id FROM events WHERE event_id = ?"
	rows, _ := db.Query(query, eventId)
	
	for rows.Next() {
		rows.Scan(&event.EventId)
	}
	
	if event.EventId != 0 {
		return true
	}
	
	return false
	
}

func GetLastEvent(db *sql.DB) *Event {
	
	event := &Event{EventId: 0}
	query := "SELECT event_id, event_name FROM events ORDER BY last_change DESC, event_id DESC LIMIT 1"
	rows, _ := db.Query(query)
	
	for rows.Next() {
		rows.Scan(&event.EventId, &event.Name)
	}
	
	if event.EventId == 0 {
		return nil
	}
	
	return event
}

func GetEventById(db *sql.DB, eventId int) *Event {
	
	event := &Event{EventId: 0}
	query := "SELECT event_id, event_name FROM events WHERE event_id = ?"
	rows, _ := db.Query(query, eventId)
	
	for rows.Next() {
		rows.Scan(&event.EventId, &event.Name)
	}
	
	if event.EventId == 0 {
		return nil
	}
	
	return event
}

func GetAllEvents(db *sql.DB) []*Event {
	
	events := []*Event{}
	
	query := "SELECT e.event_id, e.event_name, COUNT(m.match_id) FROM events e JOIN matches m ON m.event_id = e.event_id GROUP BY e.event_id ORDER BY e.last_change DESC, event_id DESC"
	
	rows, _ := db.Query(query)
	
	for rows.Next() {
		event := &Event{}
		rows.Scan(&event.EventId, &event.Name, &event.MatchesCount)
		events = append(events, event)
	}	
	
	return events
}

func ImportEvents(db *sql.DB, events []*csgoscrapper.Event) {
	
	for _, evt := range events {
		
		query := "INSERT INTO events (event_id, event_name) VALUES (?, ?)"
		
		_, err := db.Query(query, evt.EventId, evt.Name)
		

		if err != nil {
			fmt.Printf("SQL Error: %v\n", err)
		}
		
	}
	
}