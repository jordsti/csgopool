package csgopool

import (
	"time"
	"csgodb"
	"eseascrapper"
	"database/sql"
	"fmt"
)


func (w *WatcherState) FetchAndImportDay(db *sql.DB, date time.Time, all_matches []*csgodb.Match, all_teams []*csgodb.Team, all_players []*csgodb.Player, division string) ([]*csgodb.Match,[]*csgodb.Team,[]*csgodb.Player) {
	url := eseascrapper.GetScheduleURL(date, division, eseascrapper.NorthAmerica)
	
	page := url.LoadPage()
	matches := page.ParseMatches()
	
	for _, m := range matches {
		if !csgodb.IsSourceMatchIn(all_matches, csgodb.EseaSource, m.MatchId) {
			//importing matches
			w.Log.Info(fmt.Sprintf("ESEA Match [%d] not existing, retrieving it", m.MatchId))
			m.ParseMatch()
			
			//checking for teams
			team1 := csgodb.FindTeamByName(all_teams, m.Team1.Name)
			if team1 == nil {
				//team not found
				team1 = csgodb.ImportEseaTeam(db, m.Team1.Team())
				all_teams = append(all_teams, team1)
			} else if team1.EseaId == 0 {
				//need to update id
				team1.EseaId = m.Team1.TeamId
				team1.UpdateSourceId(db)
			}
			
			team2 := csgodb.FindTeamByName(all_teams, m.Team2.Name)
			if team2 == nil {
				//team not found
				team2 = csgodb.ImportEseaTeam(db, m.Team2.Team())
				all_teams = append(all_teams, team2)
			} else if team1.EseaId == 0 {
				//need to update id
				team2.EseaId = m.Team2.TeamId
				team2.UpdateSourceId(db)
			}
			
			m.Team1.TeamId = team1.TeamId
			m.Team2.TeamId = team2.TeamId
			
			//players now, must fix match stat to apply mysql id and not esea
			for _, ms := range m.PlayerStats {
				player := csgodb.FindPlayerByName(all_players, ms.Name)
				if player == nil {
					//player not found
					player = csgodb.ImportEseaPlayer(db, ms.Player())
					all_players = append(all_players, player)
				} else if player.EseaId == 0 {
					//need to update id
					player.EseaId = ms.PlayerId
					player.UpdateSourceId(db)
				}
				w.Log.Debug(fmt.Sprintf("Player Id : %d [%d]", player.PlayerId, len(all_players)))
				ms.PlayerId = player.PlayerId
				
				//fixing team id too
				if ms.TeamId == team1.EseaId {
					ms.TeamId = team1.TeamId
				} else if ms.TeamId == team2.EseaId {
					ms.TeamId = team2.TeamId
				}
				
				csgodb.AddPlayerToTeam(db, ms.TeamId, ms.PlayerId)
			}
			
			//match can be imported now !
			csgodb.ImportEseaMatch(db, m)
		}
		
	}
	
	return all_matches, all_teams, all_players
}

func (w *WatcherState) UpdateESEA(dayDelta int, division string) {
	db, _ := csgodb.Db.Open()
	
	all_matches := csgodb.GetAllMatches(db)
	all_teams := csgodb.GetAllTeams(db)
	all_players := csgodb.GetAllPlayers(db)
	
	today := time.Now()

	//fetching yesterday stats
	for it := dayDelta; it >= 0; it-- {
		date := today.AddDate(0, 0, -it)
		fmt.Printf("Checking ESEA Matches [%d-%02d-%02d]\n", date.Year(), date.Month(), date.Day())
		all_matches, all_teams, all_players = w.FetchAndImportDay(db, date, all_matches, all_teams, all_players, division)
		fmt.Printf("%d : %d", len(all_teams), len(all_players))
	}
	
	db.Close()
}