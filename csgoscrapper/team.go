package csgoscrapper

import (
	"regexp"
	"strconv"
	"fmt"
	"io/ioutil"
	"encoding/json"
)


type TeamStat struct {
	MapsPlayed int
	Wins int
	Draws int
	Losses int
	Frags int
	Deaths int
	RoundsPlayed int
	KDRatio float32
}


type Team struct {
	Name string
	TeamId int
	Stats TeamStat
	Players []Player
}
//team id must be verified before calling this !
func (t *Team) GetPlayerById(id int) *Player {
	
	for _, p := range t.Players {
		if p.PlayerId == id {
			return &p
		}
	}
	
	//need to add this player to this team
	pl := Player{PlayerId: id, Name:""}
	_pl := &pl
	_pl.LoadStats()
	t.Players = append(t.Players, pl)
	return _pl	
}

func GetTeamById(teams []*Team, id int) *Team {
	
	for _, t := range teams {
		if t.TeamId == id {
			return t
		}
	}
	
	return nil
}

func (t Team) PlayersCount() int {
	count := 0
	
	for i := range t.Players {
		if i >= 0 {
			count = count + 1
		}
	}
	
	return count
}

func SaveTeams(teams []*Team, path string) {
	b, err := json.MarshalIndent(teams, "", "	")
	
	if err != nil {
		fmt.Println("Error while saving teams [1]...")
	}
	
	err = ioutil.WriteFile(path, b, 0644)
	
	if err != nil {
		fmt.Println("Error while saving teams [2]...")
	}
}

func LoadTeams(path string) []*Team {
	teams := []*Team{}
	
	data, err := ioutil.ReadFile(path)
	
	if err != nil {
		fmt.Println("Error while reading teams [1]...")
	}
	
	err = json.Unmarshal(data, &teams)
	
	if err != nil {
		fmt.Println("Error while reading teams [2]...")
		fmt.Println(err)
	}
	
	return teams
		
}

func (t *Team) LoadTeam() {
	if t.Name == "NotSet" {
		purl := GetTeamPage(t.TeamId)
		
		page, err := purl.LoadPage()
		
		if err != nil {
			fmt.Println("Error while loading page")
		}
		
		re := regexp.MustCompile("Team stats: ([a-zA-Z0-9 \\-\\.!]+) <span")
		rs := re.FindAllStringSubmatch(page.Content, -1)
		
		t.Name = rs[0][1]
		
		//parsing initial stats
		
		//maps played
		
		re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;\">Maps played</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:180px;float:left;color:black;text-align:right;\">([0-9]+)</div>")
		rs = re.FindAllStringSubmatch(page.Content, -1)
		
		mapPlayed, _ := strconv.ParseInt(rs[0][1], 10, 32)
		t.Stats.MapsPlayed = int(mapPlayed)
		
		//wins/ draws/ losses
		
		re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:140px;float:left;\">Wins / draws / losses</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:140px;float:left;color:black;text-align:right;\">([0-9]+) / ([0-9]+) / ([0-9]+)</div>")
		rs = re.FindAllStringSubmatch(page.Content, -1)
		
		wins, _ := strconv.ParseInt(rs[0][1], 10, 32)
		draws, _ := strconv.ParseInt(rs[0][2], 10, 32)
		losses, _ := strconv.ParseInt(rs[0][3], 10, 32)
		
		t.Stats.Wins = int(wins)
		t.Stats.Draws = int(draws)
		t.Stats.Losses = int(losses)
		
		// total frags
		
		re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;\">Total kills</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:180px;float:left;color:black;text-align:right;\">([0-9]+)</div>")
		rs = re.FindAllStringSubmatch(page.Content, -1)
		
		frags, _ := strconv.ParseInt(rs[0][1], 10, 32)
		
		t.Stats.Frags = int(frags)
		
		//total deaths
		
		re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;\">Total deaths</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:180px;float:left;color:black;text-align:right;\">([0-9]+)</div>")
		rs = re.FindAllStringSubmatch(page.Content, -1)
		
		deaths, _ := strconv.ParseInt(rs[0][1], 10, 32)
		
		t.Stats.Deaths = int(deaths)
		
		re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;\">Rounds played</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:180px;float:left;color:black;text-align:right;\">([0-9]+)</div>")
		rs = re.FindAllStringSubmatch(page.Content, -1)
		
		rounds, _ := strconv.ParseInt(rs[0][1], 10, 32)
		
		t.Stats.RoundsPlayed = int(rounds)
		
		re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;\">K/D Ratio</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:180px;float:left;color:black;text-align:right;\">([0-9]+\\.?[0-9]*)</div>")
		rs = re.FindAllStringSubmatch(page.Content, -1)
		
		kdratio, _ := strconv.ParseFloat(rs[0][1], 32)
		
		t.Stats.KDRatio = float32(kdratio)
		
		t.Players = page.ParsePlayer()
		
	}
	
}

func (p PageContent) ParseTeams() []*Team {
	teams := []*Team{}
	// ?pageid=179&amp;teamid=4411&amp;eventid=0&amp;gameid=2
	
	regex := "\\?pageid=179&amp;teamid=([0-9]+)&amp;eventid="+strconv.Itoa(p.URL.EventId)+"&amp;gameid="+strconv.Itoa(p.URL.GameId)

	re := regexp.MustCompile(regex)
	
	rs := re.FindAllStringSubmatch(p.Content, -1)
	
	for _, t := range rs {
		teamId, _ := strconv.ParseInt(t[1], 10, 32)
		
		team := &Team{Name:"NotSet", TeamId: int(teamId)}
		
		teams = append(teams, team)
	}
	
	return teams
}


