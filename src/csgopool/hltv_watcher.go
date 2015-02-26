package csgopool

import (
	"csgodb"
	"csgoscrapper"
	"database/sql"
	"fmt"
)

func (w *WatcherState) FetchAndImportHltvMatches(db *sql.DB, all_matches []*csgodb.Match, all_teams []*csgodb.Team, all_players []*csgodb.Player, all_events []*csgodb.Event, pageOffset int) ([]*csgodb.Match,[]*csgodb.Team,[]*csgodb.Player, []*csgodb.Event) {
	
	matches := csgoscrapper.GetMatches(pageOffset)
	
	for _, match := range matches {
		
		if !csgodb.IsSourceMatchIn(all_matches, csgodb.HltvSource, match.MatchId) {
			match.ParseMatch()
			w.Log.Info(fmt.Sprintf("Hltv Match [%d] not found, fetching it", match.MatchId))
			
			//checking for team
			team1 := match.GetTeam1()
			team2 := match.GetTeam2()
			
			dbTeam1 := csgodb.FindTeamByName(all_teams, team1.Name)
			if dbTeam1 == nil {
				//not existing importing
				dbTeam1 = csgodb.ImportHltvTeam(db, team1)
				all_teams = append(all_teams, dbTeam1)
			} else if dbTeam1.HltvId == 0 {
				dbTeam1.HltvId = team1.TeamId
				dbTeam1.UpdateSourceId(db)
			}
			
			dbTeam2 := csgodb.FindTeamByName(all_teams, team2.Name)
			if dbTeam2 == nil {
				//not existing importing
				dbTeam2 = csgodb.ImportHltvTeam(db, team2)
				all_teams = append(all_teams, dbTeam2)
			} else if dbTeam2.HltvId == 0 {
				dbTeam2.HltvId = team2.TeamId
				dbTeam2.UpdateSourceId(db)
			}
			
			match.Team1.TeamId = dbTeam1.TeamId
			match.Team2.TeamId = dbTeam2.TeamId
			
			//players check
			for _, ms := range match.PlayerStats {
				hltvPlayer := ms.Player()
				player := csgodb.FindPlayerByName(all_players, hltvPlayer.Name)
				
				if player == nil {
					player = csgodb.ImportHltvPlayer(db, hltvPlayer)
					all_players = append(all_players, player)
				} else if player.HltvId == 0 {
					player.HltvId = hltvPlayer.PlayerId
					player.UpdateSourceId(db)
				}
				
				ms.PlayerId = player.PlayerId
				
				//need to fix team id too
				
				if ms.TeamId == team1.TeamId {
					ms.TeamId = dbTeam1.TeamId
				} else {
					ms.TeamId = dbTeam2.TeamId
				}
				
				csgodb.AddPlayerToTeam(db, ms.TeamId, ms.PlayerId)
			}
			
			_match := csgodb.ImportHltvMatch(db, match)
			all_matches = append(all_matches, _match)
		}
		
	}
	
	return all_matches, all_teams, all_players, all_events
}

func (w *WatcherState) UpdateHltv(matchesPage int) {
	
	db, _ := csgodb.Db.Open()
	
	all_matches := csgodb.GetAllMatches(db)
	all_teams := csgodb.GetAllTeams(db)
	all_players := csgodb.GetAllPlayers(db)
	all_events := csgodb.GetAllEvents(db)
	
	for it := matchesPage; it >= 0; it-- {
		pageOffset := it * 50
		
		all_matches, all_teams, all_players, all_events = w.FetchAndImportHltvMatches(db, all_matches, all_teams, all_players, all_events, pageOffset)
		
	}
	
	db.Close()
	
}