package csgodb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"io/ioutil"
	"fmt"
)


var Db *Database = &Database{}

type Database struct {
	Username string
	Password string
	Address string
	Name string
}

func (d *Database) LoadConfig(path string) {
	b, err := ioutil.ReadFile(path)
	
	if err != nil {
		fmt.Println("Error while loading database config [1]")
		return
	}
	
	err = json.Unmarshal(b, d)
	
	if err != nil {
		fmt.Println("Error while loading database config [2]")
	}
}

func (d *Database) SaveConfig(path string) {
	
	b, err := json.MarshalIndent(d, "", "	")
	
	if err != nil {
		fmt.Println("Error while saving database config [1]")
		return
	}
	
	err = ioutil.WriteFile(path, b, 0644)
	
	if err != nil {
		fmt.Println("Error while saving database config [2]")
	}
}


func (d *Database) GetDSN() string {
	
	dsn := ""
	
	if len(d.Password) > 0 {
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", d.Username, d.Password, d.Address, d.Name)
	} else {
		dsn = fmt.Sprintf("%s@%s/%s?parseTime=true", d.Username, d.Address, d.Name)
	}
	
	return dsn
}

func (d *Database) Open() (*sql.DB, error) {
	return sql.Open("mysql", d.GetDSN())
}

func InitTables(db *sql.DB) {
	//using HLTV id, no auto_increment
	//team tables
	req := "CREATE TABLE IF NOT EXISTS `teams` ("
	req += "`team_id` int(255) NOT NULL,"
	req += "`team_name` varchar(255) NOT NULL,"
	req += "PRIMARY KEY (`team_id`)"
	req +=  ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err := db.Exec(req)
	
	req = "CREATE TABLE IF NOT EXISTS `players` ("
	req += "`player_id` int(255) NOT NULL,"
	req += "`player_name` varchar(255) NOT NULL,"
	req += "PRIMARY KEY (`player_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	
	req = "CREATE TABLE IF NOT EXISTS `players_teams` ("
	req += "`player_id` int(255) NOT NULL,"
	req += "`team_id` int(255) NOT NULL"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	
	req = "CREATE TABLE IF NOT EXISTS `events` ("
	req += "`event_id` int(255) NOT NULL,"
	req += "`event_name` varchar(255) NOT NULL,"
	req += "PRIMARY KEY (`event_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	
	req = "CREATE TABLE IF NOT EXISTS `matches` ("
	req += "`match_id` int(255) NOT NULL,"
	req += "`team1_id` int(255) NOT NULL,"
	req += "`team1_score` int(255) NOT NULL,"
	req += "`team2_id` int(255) NOT NULL,"
	req += "`team2_score` int(255) NOT NULL,"
	req += "`map` varchar(255) NOT NULL,"
	req += "`event_id` int(255) NOT NULL,"
	req += "`match_date` DATE NOT NULL,"
	req += "PRIMARY KEY (`match_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	
	req = "CREATE TABLE IF NOT EXISTS `matches_stats` ("
	req += "`match_stat_id` int(255) NOT NULL AUTO_INCREMENT,"
	req += "`match_id` int(255) NOT NULL,"
	req += "`team_id` int(255) NOT NULL,"
	req += "`player_id` int(255) NOT NULL,"
	req += "`frags` int(255) NOT NULL,"
	req += "`headshots` int(255) NOT NULL,"
	req += "`assists` int(255) NOT NULL,"
	req += "`deaths` int(255) NOT NULL,"
	req += "`kdratio` float NOT NULL,"
	req += "`kddelta` int(25) NOT NULL,"
	req += "PRIMARY KEY (`match_stat_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	
	req = "CREATE TABLE IF NOT EXISTS `users` ("
	req += "`user_id` int(255) NOT NULL AUTO_INCREMENT,"
	req += "`username` varchar(255) NOT NULL,"
	req += "`password` varchar(255) NOT NULL,"
	req += "`email` varchar(255) NOT NULL,"
	req += "`rank`int(25) NOT NULL,"
	req += "PRIMARY KEY(`user_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
}