package csgodb

import (
	"database/sql"
	"steamapi"
)

type SteamID struct {
	LinkId int
	SteamId int64
	UserId int
}

func GetUserBySteamID(db *sql.DB, steamId int64) *User {
	
	query := "SELECT user_id FROM steams_id WHERE steam_id = ?"
	rows, _ := db.Query(query, steamId)
	for rows.Next() {
		uid := 0
		rows.Scan(&uid)
		
		user := GetUserById(db, uid)
		if user.Id != 0 {
			return user
		}
	}
	
	return nil
} 

func (sid *SteamID) ApiID() *steamapi.SteamID {
	
	steamid := &steamapi.SteamID{}
	steamid.AccountId = sid.SteamId
	steamid.Type = steamapi.IndividualType
	
	return steamid
}

func GetSteamIDByUser(db *sql.DB, userId int) *SteamID {
	
	steamid := &SteamID{}
	
	query := "SELECT link_id, steam_id, user_id FROM steams_id WHERE user_id = ?"
	
	rows, _ := db.Query(query, userId)
	for rows.Next() {
		rows.Scan(&steamid.LinkId, &steamid.SteamId, &steamid.UserId)
	}
	
	return steamid
}

func (sid *SteamID) UpdateSteamID(db *sql.DB) {
	
	query := "UPDATE steams_id SET steam_id = ? WHERE user_id = ?"
	db.Exec(query,sid.SteamId, sid.UserId)
	
}

func SaveSteamID(db *sql.DB, userId int, steamId int64) {
	query := "INSERT INTO steams_id (steam_id, user_id) VALUES (?, ?)"
	db.Exec(query, steamId, userId)
}