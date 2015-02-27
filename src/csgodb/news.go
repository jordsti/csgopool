package csgodb

import (
	"time"
	"database/sql"
)

type Author struct {
	AuthorId int
	Name string
}

type News struct {
	NewsId int
	Title string
	Text string
	Author Author
	PostedOn time.Time
}

func (n *News) Insert(db *sql.DB) {
	time := time.Now()
	query := `INSERT INTO news (news_title, news_text, author_id, posted_on) VALUES (?, ?, ?, ?)`
	db.Exec(query, n.Title, n.Text, n.Author.AuthorId, time)
}

func GetLastNews(db *sql.DB) *News {
	news := &News{}
	
	query := `SELECT n.news_id, n.news_title, n.news_text, n.author_id, u.username, n.posted_on 
		FROM news n 
		JOIN users u ON u.user_id = n.author_id 
		ORDER BY news_id DESC LIMIT 1`
	
	rows, _ := db.Query(query)
	for rows.Next() {
		rows.Scan(&news.NewsId, &news.Title, &news.Text, &news.Author.AuthorId, &news.Author.Name, &news.PostedOn)
	}
	
	return news
}