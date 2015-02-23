package main

import (
	"csgodb"
	"fmt"
)

type StackInfo struct {
	GitUrl string
	SnapshotUrl string
	Port int
	DataPath string
	WebRoot string
	Id string
	DbPassword string
	Dependencies []string
	Env *Environment
}

func DefaultPoolStack() *StackInfo {
	stack := &StackInfo{}
	
	stack.GitUrl = "https://github.com/jordsti/csgopool.git"
	stack.SnapshotUrl = "http://csgopool.com/snapshots/snapshot-current.json"
	stack.Port = 8000
	
	stack.Dependencies = append(stack.Dependencies, "golang.org/x/crypto/bcrypt")
	stack.Dependencies = append(stack.Dependencies, "github.com/go-sql-driver/mysql")
	
	return stack
}


func (s *StackInfo) GenerateId() {
	s.Id = csgodb.RandomString(10)
	s.DbPassword = csgodb.RandomString(18)
	
}

func (s *StackInfo) PrepareDatabase() {
	
	conn, err := db.Open()
	if err != nil {
		fmt.Printf("%v\n", err)
		
	}
	db_user := "s_" + s.Id
	
	query := fmt.Sprintf("create database %s", db_user)
	_, err = conn.Exec(query)
	if err != nil {
		fmt.Printf("%v\n", err)
		
	}
	query = fmt.Sprintf("grant usage on *.* to %s@localhost identified by '%s'", db_user, s.DbPassword)
	_, err = conn.Exec(query)
		if err != nil {
		fmt.Printf("%v\n", err)
		
	}
	query = fmt.Sprintf("grant all privileges on %s.* to %s@localhost", db_user, db_user)
	conn.Exec(query)
	
	conn.Close()
	
	
}