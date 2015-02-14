package csgoscrapper

import (
	"net/http"
	"strconv"
	"io/ioutil"
)

type PageURL struct {
	Name string
	PageId int
	EventId int
	GameId int
	TeamId int
	PlayerId int
	MatchId int
}

type PageContent struct {
	URL PageURL
	Content string
	Status int
}

func GetMatchPage(matchId int) PageURL {
	p := PageURL {"MatchPage", 188, -1, 2, -1, -1, matchId}
	return p
}

func GetTeamsPage() PageURL {
	p := PageURL{"TeamsPage", 182, 0, 2, -1, -1, -1}
	return p
}

func GetTeamPage(teamId int) PageURL {
	p := PageURL{"TeamPage", 179, 0, 2, teamId, -1, -1}
	return p
}

func GetPlayerPage(playerId int) PageURL {
	p := PageURL{"PlayerPage", 173, 0, 2, -1, playerId, -1}
	return p
}

func GetEventsPage() PageURL {
	p := PageURL{"EventsPage", 193, -1, 2, -1, -1, -1}
	return p
}

func GetEventMatches(eventId int) PageURL {
	p := PageURL{"EventMatchesPage", 188, eventId, 2, -1, -1, -1}
	return p
}

func (p PageURL) GenerateURL() string {
	base_url := "http://www.hltv.org/?"
	base_url = base_url + "pageid=" + strconv.Itoa(p.PageId) + "&gameid=" + strconv.Itoa(p.GameId)
	
	if p.EventId != -1 {
		base_url = base_url + "&eventid=" + strconv.Itoa(p.EventId)
	}
	
	if p.TeamId != -1 {
		base_url = base_url + "&teamid=" + strconv.Itoa(p.TeamId)
	}
	
	if p.PlayerId != -1 {
		base_url = base_url + "&playerid=" + strconv.Itoa(p.PlayerId)
	}
	
	if p.MatchId != -1 {
		base_url = base_url + "&matchid=" + strconv.Itoa(p.MatchId)
	}
	
	return base_url
}

func (p PageURL) LoadPage() (PageContent, error) {
	//building url
	base_url := p.GenerateURL()
	resp, err := http.Get(base_url)
	
	defer resp.Body.Close()
	
	if err != nil {
		return PageContent{}, err
	} else {
		body, err2 := ioutil.ReadAll(resp.Body)
		
		if err2 != nil {
			return PageContent{}, err2
		}
		
		
		content := PageContent{URL: p, Content: string(body), Status: resp.StatusCode}
		
		return content, nil
	}
	
}
