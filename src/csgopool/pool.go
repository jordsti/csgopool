package csgopool

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"csgodb"
	"database/sql"
)


type PoolSetting struct {
	PoolOn bool
	AutoAddMatches bool
}


type PoolState struct {
	Settings PoolSetting
	Path string
}


var Pool *PoolState

func NewPoolState(settingPath string) {
	
	Pool = &PoolState{}
	Pool.Settings.PoolOn = false
	Pool.Settings.AutoAddMatches = false
	
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

func AttributePoints(db *sql.DB, matchId int) {
	
	match := csgodb.GetMatchById(db, matchId)
	if match != nil {
		match.FetchStats(db)
		
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
			
			points = (ps.Frags * 2 + ps.Headshots * 2 + ps.Assists/2) - ps.Deaths
			
			//kdddelta bonus if positive
			kddelta := ps.KDDelta * 3
			
			if kddelta > 0 {
				points += kddelta
			}
			
			if points < 0 {
				points = 0
			}
			
			csgodb.AddPoint(db, match.MatchId, ps.PlayerId, points)
		}
	}
	
}