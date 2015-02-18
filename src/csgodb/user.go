package csgodb

import (
	"errors"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"encoding/hex"
	"crypto/rand"
)

const (
	UserRank = 1
	PoolMaster = 10
)

type UserConstraint struct {
	NameMax int
	NameMin int
	PasswordMin int
}

type User struct {
	Id int
	Name string
	Password string
	Email string
	Rank int
}

type Users struct {
	Users []User
	CurrentId int
	Constraints UserConstraint
}

func (u User) IsPoolMaster() bool {
	if u.Rank >= PoolMaster {
		return true
	}
	
	return false
}

func UsersCount(db *sql.DB) int {
	query := "SELECT user_id FROM users"
	rows, _ := db.Query(query)
	count := 0
	for rows.Next() {
		count++
	}
	
	return count
}

func Login(db *sql.DB, username string, password string) (*User, error) {
	
	user := GetUserByName(db, username)
	
	if user == nil {
		return nil, errors.New("Invalid username and/or password")
	}
	
	b, _ := hex.DecodeString(user.Password)
	
	err := bcrypt.CompareHashAndPassword(b, []byte(password))
	
	if err == nil {
		//success
		return user, nil
	} else {
		return nil, errors.New("Invalid username and/or password")
	}
	
}

func CreateUser(db *sql.DB, username string, password string, email string, rank int) error {
	
	//username unique check
	user := GetUserByName(db, username)
	
	if user != nil {
		return errors.New("Username already exists")
	}
	
	user = GetUserByEmail(db, email)
	
	if user != nil {
		return errors.New("Email already exists")
	}
	
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	pw_str := hex.EncodeToString(pwd)
	
	nuser := User{Name: username, Password: pw_str, Email:email, Rank:rank}
	
	query := "INSERT INTO users (username, password, email, rank) VALUES (?, ?, ?, ?)"
	db.Exec(query, nuser.Name, nuser.Password, nuser.Email, nuser.Rank)
	
	return nil
}

func GetUserByName(db *sql.DB, username string) *User {
	
	user := &User{Id: 0}
	
	query := "SELECT user_id, username, password, email, rank FROM users WHERE username = ?"
	rows, _ := db.Query(query, username)
	
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Rank)
	}
	
	if user.Id != 0 {
		return user
	} else {
		return nil
	}
}


func GetUserByEmail(db *sql.DB, email string) *User {

	user := &User{Id: 0}
	
	query := "SELECT user_id, username, password, email, rank FROM users WHERE email = ?"
	rows, _ := db.Query(query, email)
	
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Rank)
	}
	
	if user.Id != 0 {
		return user
	} else {
		return nil
	}
}

func RandomString(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0987654321"
	
	var bytes = make([]byte, length)
	rand.Read(bytes)
	
	for k, v := range bytes {
		bytes[k] = chars[v%byte(len(chars))]
	}
	
	return string(bytes)
	
}
