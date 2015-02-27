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
	Location string
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
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=%s", d.Username, d.Password, d.Address, d.Name, d.Location)
	} else {
		dsn = fmt.Sprintf("%s@tcp(%s)/%s?parseTime=true&loc=%s", d.Username, d.Address, d.Name, d.Location)
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
	req += "`team_id` int(255) NOT NULL AUTO_INCREMENT,"
	req += "`team_name` varchar(255) NOT NULL,"
	req += "`hltv_id` int(255) NOT NULL,"
	req += "`esea_id` int(255) NOT NULL,"
	req += "UNIQUE KEY `teams_constraint` (`team_name`),"
	req += "PRIMARY KEY (`team_id`)"
	req +=  ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err := db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `players` ("
	req += "`player_id` int(255) NOT NULL AUTO_INCREMENT,"
	req += "`player_name` varchar(255) NOT NULL,"
	req += "`player_alias` varchar(255) NOT NULL,"
	req += "`hltv_id` int(255) NOT NULL,"
	req += "`esea_id` int(255) NOT NULL,"
	req += "UNIQUE KEY `players_constraint` (`player_name`),"
	req += "PRIMARY KEY (`player_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `players_teams` ("
	req += "`link_id` int(255) NOT NULL AUTO_INCREMENT,"
	req += "`player_id` int(255) NOT NULL,"
	req += "`team_id` int(255) NOT NULL,"
	req += "UNIQUE KEY `players_teams_constraint` (`player_id`, `team_id`),"
	req += "PRIMARY KEY (`link_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `events` ("
	req += "`event_id` int(255) NOT NULL AUTO_INCREMENT,"
	req += "`source_id` int(255) NOT NULL, "
	req += "`source` int(25) NOT NULL,"
	req += "`event_name` varchar(255) NOT NULL,"
	req += "`last_change` DATETIME NOT NULL,"
	req += "PRIMARY KEY (`event_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `matches` ("
	req += "`match_id` int(255) NOT NULL AUTO_INCREMENT,"
	req += "`team1_id` int(255) NOT NULL,"
	req += "`team1_score` int(255) NOT NULL,"
	req += "`team2_id` int(255) NOT NULL,"
	req += "`team2_score` int(255) NOT NULL,"
	req += "`map` varchar(255) NOT NULL,"
	req += "`event_id` int(255) NOT NULL,"
	req += "`match_date` DATE NOT NULL,"
	req += "`source_id` int(255) NOT NULL,"
	req += "`source` int(25) NOT NULL,"
	req += "`pool_status` int(255) NOT NULL, "
	req += "UNIQUE KEY `matches_constraint` (`source_id`, `source`), "
	req += "PRIMARY KEY (`match_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `matches_stats` ("
	req += "`match_stat_id` int(255) NOT NULL AUTO_INCREMENT,"
	req += "`match_id` int(255) NOT NULL,"
	req += "`team_id` int(255) NOT NULL,"
	req += "`player_id` int(255) NOT NULL,"
	req += "`frags` int(255) NOT NULL,"
	//req += "`headshots` int(255) NOT NULL," removed because ESEA doesnt show that stats!
	req += "`assists` int(255) NOT NULL,"
	req += "`deaths` int(255) NOT NULL,"
	req += "`kdratio` float NOT NULL,"
	req += "`kddelta` int(25) NOT NULL,"
	req += "UNIQUE KEY `matches_stats_constraint` (`match_id`, `player_id`), "
	req += "PRIMARY KEY (`match_stat_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
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
	req = "CREATE TABLE IF NOT EXISTS `divisions` ("
	req += "`division_id` int(255) NOT NULL AUTO_INCREMENT, "
	req += "`division_name` varchar(255) NOT NULL, "
	req += "PRIMARY KEY(`division_id`) "
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `divisions_players` ( "
	req += "`link_id` int(255) NOT NULL AUTO_INCREMENT, "
	req += "`division_id` int(255) NOT NULL, "
	req += "`player_id` int(255) NOT NULL,"
	req += "UNIQUE KEY `divisions_players_constraint` (`division_id`, `player_id`),"
	req += "PRIMARY KEY (`link_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `users_pools` ( "
	req += "`pool_id` int(255) NOT NULL AUTO_INCREMENT, "
	req += "`division_id` int(255) NOT NULL, "
	req += "`user_id` int(255) NOT NULL, "
	req += "`player_id` int(255) NOT NULL, "
	req += "`created_on` DATETIME NOT NULL, "
	req += "UNIQUE KEY `users_pools_constraint` (`division_id`, `user_id`),"
	req += "PRIMARY KEY (`pool_id`)"
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `watcher_update` ( "
	req += "`update_id` int(255) NOT NULL AUTO_INCREMENT, "
	req += "`update_time` DATETIME NOT NULL, "
	req += "PRIMARY KEY(`update_id`) "
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `players_points` ( "
	req += "`point_id` int(255) NOT NULL AUTO_INCREMENT, "
	req += "`player_id` int(255) NOT NULL, "
	req += "`match_id` int(255) NOT NULL, "
	req += "`points` int(255) NOT NULL, "
	req += "UNIQUE KEY `players_points_constraint` (`player_id`, `match_id`),"
	req += "PRIMARY KEY(`point_id`) " 
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	req = "CREATE TABLE IF NOT EXISTS `monitor_tasks` ( "
	req += "`task_id` int(255) NOT NULL AUTO_INCREMENT, "
	req += "`task_type` int(25) NOT NULL, "
	req += "`task_status` int(25) NOT NULL, "
	req += "`task_date` DATETIME NOT NULL, "
	req += "PRIMARY KEY(`task_id`) "
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	
	req = "CREATE TABLE IF NOT EXISTS `news` ( "
	req += "`news_id` int(255) NOT NULL AUTO_INCREMENT, "
	req += "`news_title` varchar(255) NOT NULL, "
	req += "`news_text` text NOT NULL, "
	req += "`author_id` int(255) NOT NULL, "
	req += "`posted_on` DATETIME NOT NULL, "
	req += "PRIMARY KEY(`news_id`) "
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
	
	req = "CREATE TABLE IF NOT EXISTS `sources` ("
	req += "`source_id` int(255) NOT NULL, "
	req += "`source_name` varchar(255) NOT NULL, "
	req += "PRIMARY KEY(`source_id`) "
	req += ") ENGINE=InnoDB CHARSET=latin1;"
	_, err = db.Exec(req)
	if err != nil {
		fmt.Printf("Database Error %v\n", err)
	}
}