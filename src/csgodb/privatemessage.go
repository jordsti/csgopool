package csgodb

import (
	"database/sql"
	"time"
)

const (
	UnreadStatus = 1
	ReadedStatus = 2
)

type PrivateMessage struct {
	MessageId int
	SenderId int
	SenderName string
	RecipientId int
	RecipientName string
	Status int
	Title string
	Text string
	SendedOn time.Time
}

func GetMessageById(db *sql.DB, messageId int) *PrivateMessage {
	msg := &PrivateMessage{}
	query := `SELECT m.message_id, m.sender_id, us.username, m.recipient_id, ur.username, m.message_status, m.message_title, m.message_text, m.sended_on
	FROM private_messages m
	JOIN users us ON us.user_id = m.sender_id
	JOIN users ur ON ur.user_id = m.recipient_id
	WHERE m.message_id = ?`
	
	rows, _ := db.Query(query, messageId)
	for rows.Next() {
		rows.Scan(&msg.MessageId, &msg.SenderId, &msg.SenderName, &msg.RecipientId, &msg.RecipientName, &msg.Status, &msg.Title, &msg.Text, &msg.SendedOn)
	}
	
	return msg
}

//title and text must be html escaped!
func AddMessage(db *sql.DB, senderId int, recipientId int, title string, text string, status int) {
	now := time.Now()
	query := `INSERT INTO private_messages (sender_id, recipient_id, message_title, message_text, message_status, sended_on) VALUES (?, ?, ?, ?, ?, ?)`
	db.Exec(query, senderId, recipientId, title, text, status, now)
}

func GetUserMessages(db *sql.DB, userId int, start int, end int) []*PrivateMessage {
	messages := []*PrivateMessage{}
	query := `SELECT m.message_id, m.sender_id, us.username, m.recipient_id, ur.username, m.message_status, m.message_title, m.message_text, m.sended_on
	FROM private_messages m
	JOIN users us ON us.user_id = m.sender_id
	JOIN users ur ON ur.user_id = m.recipient_id
	WHERE m.recipient_id = ? 
	ORDER BY m.sended_on DESC
	LIMIT ?, ?`
	
	rows, _ := db.Query(query, userId, start, end)
	for rows.Next() {
		msg := &PrivateMessage{}
		rows.Scan(&msg.MessageId, &msg.SenderId, &msg.SenderName, &msg.RecipientId, &msg.RecipientName, &msg.Status, &msg.Title, &msg.Text, &msg.SendedOn)
		
		messages = append(messages, msg)
	}
	
	return messages
}

func (pm *PrivateMessage) UpdateStatus(db *sql.DB, status int) {
	query := `UPDATE private_messages SET message_status = ? WHERE message_id = ?`
	db.Exec(query, status, pm.MessageId)
}