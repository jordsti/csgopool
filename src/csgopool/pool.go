package csgopool

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"csgodb"
	"database/sql"
)

type MailSetting struct {
	Address string
	Port int
	Username string
	Password string
	Email string
}

type PoolSetting struct {
	PoolOn bool
	AutoAddMatches bool
	SteamKey string
	SteamBot bool
	PoolCost float32
	MailVerification bool
	Mail MailSetting
}


type PoolState struct {
	Settings PoolSetting
	Path string
}


var Pool *PoolState

func (ms MailSetting) Host() string {
	return fmt.Sprintf("%s:%d", ms.Address, ms.Port)
}

func NewPoolState(settingPath string) {
	
	Pool = &PoolState{}
	Pool.Settings.PoolOn = false
	Pool.Settings.AutoAddMatches = false
	Pool.Settings.SteamBot = false
	Pool.Settings.PoolCost = float32(0.00)
	Pool.Settings.MailVerification = false
	Pool.LoadSetting(settingPath)
	
} 

func (p *PoolState) SaveSetting(path string) {
	
	b, err := json.MarshalIndent(p.Settings, "", "	")
	
	if err != nil {
		watcher.Log.Error("JSON encoding of setting faile")
	}
	
	err = ioutil.WriteFile(path, []byte(b), 0644)
	
	if err != nil {
		watcher.Log.Error("Error while writing setting file")
	}
	
	p.Path = path
}

func (p *PoolState) LoadSetting(path string) {
	
	b, err := ioutil.ReadFile(path)
	if err != nil {
		//writing file
		watcher.Log.Info(fmt.Sprintf("Setting file not found creating default settings at %s", path))
		p.SaveSetting(path)
		return
	}
	
	err = json.Unmarshal(b, &p.Settings)
	
	if err != nil {
		watcher.Log.Error("Error while parsing settings")
	}
	
	p.Path = path
	
}

func RemovePoints(db *sql.DB, matchId int) {
	match := csgodb.GetMatchById(db, matchId)
	
	if match != nil {
		
		if match.PoolStatus == 1 {
			
			query := "DELETE FROM players_points WHERE match_id = ?"
			db.Exec(query, matchId)
			
		}
		
		csgodb.UpdateMatchPoolStatus(db, matchId, 0)
	}
}

func AttributePoints(db *sql.DB, matchId int) {
	
	match := csgodb.GetMatchById(db, matchId)
	
	if match != nil {
		match.FetchStats(db)
		
		winId := 0
		tie := false
		
		if match.Team1.Score == match.Team2.Score {
			//tie
			tie = true
		} else if match.Team1.Score > match.Team2.Score {
			//team1 win
			winId = match.Team1.TeamId
		} else {
			winId = match.Team2.TeamId
		}
		
		for _, ps := range match.PlayerStats {
			/*
				type PlayerStat struct {
					MatchStatId int
					MatchId int
					TeamId int
					PlayerId int
					Frags int
					Headshots int
					Assists int
					Deaths int
					KDRatio float32
					KDDelta int
					PlayerName string
				}
			*/
			
			points := 0
			
			points += (ps.Frags * 2 + ps.Assists)
			
			//kdddelta bonus if positive
			kddelta := ps.KDDelta * 3
			
			if kddelta > 0 {
				points += kddelta
			}
			
			if points < 0 {
				points = 0
			}
			
			//win or tie points awards
			if tie {
				points += 10
			} else if ps.TeamId == winId {
				points += 25
			}
			
			csgodb.AddPoint(db, match.MatchId, ps.PlayerId, points)
		}
		
		csgodb.UpdateMatchPoolStatus(db, matchId, 1)
	}
	
}