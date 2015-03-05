package csgodb

import (
  "database/sql"
  "time"
  "encoding/json"
  "fmt"
)

type ItemInfo struct {
    HashName string
    ClassId string
    LocalId uint64
}

type TransactionData struct {
    Item *ItemInfo
}

type Transaction struct {
    TransactionId int
    UserId int
    Timestamp time.Time
    Description string
    Data string //maybe json data
    Amount float32
}

func (ts *Transaction) SetData(td *TransactionData) {
 
  b, err := json.Marshal(td)
  if err != nil {
    fmt.Printf("JSON Error : %v\n", err)
  }
  
  ts.Data = string(b)
}

func (ts *Transaction) GetData() *TransactionData {
  td := &TransactionData{} 
  
  err := json.Unmarshal([]byte(ts.Data), td)
  if err != nil {
    fmt.Printf("JSON Error : %v\n", err)
  }
  
  return td
}

func GetTransactionsByUser(db *sql.DB, userId int) []*Transaction {
    transactions := []*Transaction{}
    
    query := "SELECT transaction_id, user_id, timestamp, description, data, amount FROM transactions WHERE user_id = ? ORDER BY transaction_id DESC"
    
    rows, _ := db.Query(query, userId)
    for rows.Next() {
		ts := &Transaction{}
		rows.Scan(&ts.TransactionId, &ts.UserId, &ts.Timestamp, &ts.Description, &ts.Data, &ts.Amount)
		transactions = append(transactions, ts)
    }
    
    return transactions
}

func (ts *Transaction) Insert(db *sql.DB) {
     
    if ts.UserId != 0 && ts.TransactionId == 0 {
        now := time.Now()
	query := "INSERT INTO transactions (user_id, timestamp, description, data, amount) VALUES (?, ?, ?, ?, ?)"
	db.Exec(query, ts.UserId, now, ts.Description, ts.Data, ts.Amount)
    }
  
}