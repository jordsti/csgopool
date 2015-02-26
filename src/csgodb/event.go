package csgodb
import (
	"database/sql"
	"hltvscrapper"
	"fmt"
	"time"
)

type Event struct {
	EventId int
	Name string
	Source int
	SourceId int
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

func FindEventByName(events []*Event, name string) *Event {
	
	for _, evt := range events {
		if evt.Name == name {
			return evt
		}
	}
	
	return nil
}

func GetAllEvents(db *sql.DB) []*Event {
	
	events := []*Event{}
	
	query := "SELECT e.event_id, e.source, e.source_id, e.event_name, COUNT(m.match_id) FROM events e LEFT JOIN matches m ON m.event_id = e.event_id GROUP BY e.event_id ORDER BY e.last_change DESC, event_id DESC"
	
	rows, _ := db.Query(query)
	
	for rows.Next() {
		event := &Event{}
		rows.Scan(&event.EventId, &event.Source, &event.SourceId, &event.Name, &event.MatchesCount)
		events = append(events, event)
	}	
	
	return events
}

func ImportHltvEvent(db *sql.DB, event *hltvscrapper.Event) *Event {
	
	query := "INSERT INTO events (source, source_id, event_name) VALUES (?, ?, ?)"
	db.Exec(query, HltvSource, event.EventId, event.Name)
	
	return GetEventBySource(db, HltvSource, event.EventId)
}

func GetEventBySource(db *sql.DB, source int, sourceId int) *Event {
	evt := &Event{}
	query := "SELECT event_id, event_name, source, source_id, last_change FROM events WHERE source = ? AND source_id = ?"
	
	rows, _ := db.Query(query, source, sourceId)
	for rows.Next() {
		rows.Scan(&evt.EventId, &evt.Name, &evt.Source, &evt.SourceId, &evt.LastChange)
	}
	
	return evt
}

func ImportEvents(db *sql.DB, events []*hltvscrapper.Event) {
	
	for _, evt := range events {
		
		query := "INSERT INTO events (event_id, event_name) VALUES (?, ?)"
		
		_, err := db.Query(query, evt.EventId, evt.Name)
		

		if err != nil {
			fmt.Printf("SQL Error: %v\n", err)
		}
		
	}
	
}