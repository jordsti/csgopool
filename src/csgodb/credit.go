package csgodb

import (
	"database/sql"
)

type Credit struct {
	CreditId int32
	UserId int32
	Amount float32
}

func InsertCredit(db *sql.DB, userId int, amount float32) {
	query := "INSERT INTO users_credit (user_id, credit) VALUES (?, ?)"
	db.Exec(query, userId, amount)	
}

func GetCreditByUser(db *sql.DB, userId int) *Credit {
	credit := &Credit{CreditId: 0}
	
	query := "SELECT credit_id, user_id, credit FROM users_credit WHERE user_id = ?"
	rows, _ := db.Query(query, userId)
	for rows.Next() {
		rows.Scan(&credit.CreditId, &credit.UserId, &credit.Amount)
	}
	
	return credit
}

func (c *Credit) Add(amount float32) {
	c.Amount += amount
}

func (c *Credit) Substract(amount float32) {
	c.Amount -= amount
}

func (c *Credit) UpdateCredit(db *sql.DB) {
	query := "UPDATE users_credit SET credit = ? WHERE user_id = ?"
	db.Exec(query, c.Amount, c.UserId)
}