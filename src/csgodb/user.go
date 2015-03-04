package csgodb

import (
	"errors"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"encoding/hex"
	"crypto/rand"
	"regexp"
	"fmt"
)

const (
	UserRank = 1
	PoolMaster = 10
)

type User struct {
	Id int
	Name string
	Password string
	Email string
	Rank int
}

type UserConstraint struct {
	NameMax int
	NameMin int
	PasswordMin int
}

func DefaultUserConstraints() *UserConstraint {
	constraints := &UserConstraint{NameMax: 12, NameMin: 4, PasswordMin: 8}
	return constraints
	
}

func (u *User) IsPoolMaster() bool {
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

func GetAllUsers(db *sql.DB) []*User {
	users := []*User{}
	query := "SELECT user_id, username, email, rank FROM users ORDER BY user_id"
	
	rows, _ := db.Query(query)
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Rank)
		users = append(users, user)
	}
	
	return users
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

func UpdatePassword(db *sql.DB, userId int, password string) error {
	constraints := DefaultUserConstraints()
	re := regexp.MustCompile(fmt.Sprintf(`^.{%d,}$`, constraints.PasswordMin))
	
	//password length
	if !re.MatchString(password) {
		return errors.New(fmt.Sprintf("Your password must contains at least %d characters", constraints.PasswordMin))
	}
	
	query := "UPDATE users SET password = ? WHERE user_id = ?"
	
	pwd, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	pw_str := hex.EncodeToString(pwd)
	
	db.Exec(query, pw_str, userId)
	
	return nil
}

func GetUserById(db *sql.DB, userId int) *User {
	user := &User{}
	
	query := "SELECT user_id, username, email, rank FROM users WHERE user_id = ?"
	rows, _ := db.Query(query, userId)
	
	for rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Rank)
	}
	
	return user
}

func CreateUser(db *sql.DB, username string, password string, email string, rank int) error {
	
	constraints := DefaultUserConstraints()
	
	//username regexp check
	re := regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9\-_]{%d,%d}$`, constraints.NameMin, constraints.NameMax))
	
	if !re.MatchString(username) {
		return errors.New(fmt.Sprintf("Username must only contains alpha-numeric characters and be between %d and %d characters", constraints.NameMin, constraints.NameMax))
	}
	
	re = regexp.MustCompile(fmt.Sprintf(`^.{%d,}$`, constraints.PasswordMin))
	
	//password length
	if !re.MatchString(password) {
		return errors.New(fmt.Sprintf("Your password must contains at least %d characters", constraints.PasswordMin))
	}
	
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
