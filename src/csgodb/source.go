package csgodb

import (
	"database/sql"
)

const (
	HltvSource = 1
	EseaSource = 2
)

type Source struct {
	SourceId int
	Name string
}

func GetAllSources(db *sql.DB) []*Source {
	sources := []*Source{}
	
	query := "SELECT source_id, source_name FROM sources"
	rows, _ := db.Query(query)
	for rows.Next() {
		source := &Source{}
		rows.Scan(&source.SourceId, &source.Name)
		sources = append(sources, source)
	}
	
	return sources
}

func InsertSource(db *sql.DB, sourceId int, name string) {
	query := "INSERT INTO sources (source_id, source_name) VALUES (?, ?)"
	db.Exec(query, sourceId, name)
}

func FindSourceById(sources []*Source, sourceId int) *Source {
	
	for _, src := range sources {
		if src.SourceId == sourceId {
			return src
		}
	}
	
	return nil
}