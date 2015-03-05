package hltvscrapper

import (
	"regexp"
	"strconv"
	"fmt"
	"github.com/moovweb/gokogiri"
	"strings"
)

type MatchTeam struct {
	TeamId int
	Name string
	Score int
}

type MatchPlayerStat struct {
	PlayerId int
	PlayerName string
	TeamId int
	TeamName string
	Frags int
	Headshots int
	Assists int
	Deaths int
	KDRatio float32
	KDDelta int
	Rating float32
}

type MatchDate struct {
	Day int
	Month int
	Year int
}

func (md MatchDate) String() string {
	return fmt.Sprintf("%02d/%02d %d", md.Day, md.Month, md.Year)
}

type MatchEvent struct {
	EventId int
	Name string
}

type Match struct {
	MatchId int
	Date MatchDate
	Team1 MatchTeam
	Team2 MatchTeam
	Map string
	Event MatchEvent
	PlayerStats []*MatchPlayerStat
}

func (mps *MatchPlayerStat) Player() Player {
	return Player{Name: mps.PlayerName, PlayerId: mps.PlayerId}
}

func (m *Match) GetTeam1() *Team {
	team := &Team{TeamId: m.Team1.TeamId, Name: m.Team1.Name}
	return team
}

func (m *Match) GetTeam2() *Team {
	
	team := &Team{TeamId: m.Team2.TeamId, Name: m.Team2.Name}
	return team
}

func (m *Match) ParseMatch() {
	
	//fmt.Printf("Parse match [%d] from hltv.org\n", m.MatchId)
	
	url := GetMatchPage(m.MatchId)
	pc, _ := url.LoadPage()
	//error handling
	
	for pc.Status != 200 {
		log.Error(fmt.Sprintf("Page Match [%d] retrieve failed, retrying now ...", m.MatchId))
		pc, _ = url.LoadPage()
	}
	
	doc, _ := gokogiri.ParseHtml([]byte(pc.Content))
	
	//parse map
	nodes, _ := doc.Search("//div[@class='covSmallHeadline' and @style='font-weight:normal;width:180px;float:left;text-align:right;']")
	//first node is the map
	m.Map = strings.TrimSpace(nodes[0].Content())
	
	//player and stats
	
	nodes, _ = doc.Search("//div[starts-with(@style,'width:606px;height:22px;background-color:')]/div[@style='padding-left:5px;padding-top:5px;']")
	
	rePlayerId := regexp.MustCompile(`playerid=([0-9]+)`)
	reTeamId := regexp.MustCompile(`teamid=([0-9]+)`)
	rePlayerScore := regexp.MustCompile(`([0-9\-]+) \(([0-9\-]+)\)`)
	
	for _, n := range nodes {
		//fmt.Printf("Id: %d\n%s\n", i, n.String())
		//player name and id
		pnames, _ := n.Search("./div[@class='covSmallHeadline' and @style='font-weight:normal;width:20%;float:left;']/a")
		
		rs := rePlayerId.FindStringSubmatch(pnames[0].Attribute("href").Value())
		
		playerName := strings.TrimSpace(pnames[0].Content())
		_playerId, _ := strconv.ParseInt(rs[1], 10, 32)
		
		//team name and id
		tnames, _ := n.Search("./div[@class='covSmallHeadline' and @style='font-weight:normal;width:20%;float:left;text-align:center']/a")
		
		rs = reTeamId.FindStringSubmatch(tnames[0].Attribute("href").Value())
		teamName := strings.TrimSpace(tnames[0].Content())
		_teamId, _ := strconv.ParseInt(rs[1], 10, 32)
		
		//stats
		stats, _ := n.Search("./div[@class='covSmallHeadline' and @style='font-weight:normal;width:10%;float:left;text-align:center']")
		
		ms := &MatchPlayerStat{}
		
		ms.PlayerId = int(_playerId)
		ms.PlayerName = playerName
		ms.TeamId = int(_teamId)
		ms.TeamName = teamName
		
		rs = rePlayerScore.FindStringSubmatch(stats[0].Content())
		//fmt.Printf("%v\n%s\n", rs, stats[0].String())
		frags, _ := strconv.ParseInt(rs[1], 10, 32)
		if rs[1] == "-" {
			frags = 0
		}
		
		assists := int64(0)
		
		if stats[1].Content() != "-" {
			assists, _ = strconv.ParseInt(stats[1].Content(), 10, 32)
		}
		
		deaths := int64(0)
		
		if stats[2].Content() != "-" {
			deaths, _ = strconv.ParseInt(stats[2].Content(), 10, 32)
		}

		ms.Frags = int(frags)
		ms.Assists = int(assists)
		ms.Deaths = int(deaths)
		ms.KDDelta = ms.Frags - ms.Deaths
		
		if deaths > 0 {
			ms.KDRatio = float32(frags) / float32(deaths)
		} else {
			ms.KDRatio = float32(frags)
		}
		
		m.PlayerStats = append(m.PlayerStats, ms)
	}	
}

func GetMatches(offset int) []*Match {
	matches := []*Match{}
	
	url := GetMatchesPage(offset)
	//error handling
	pc, _ := url.LoadPage()
	//fmt.Printf(pc.Content)
	for pc.Status != 200 {
		log.Error(fmt.Sprintf("Page could not be retrieve : %s, retrying", url.GenerateURL()))
		pc, _ = url.LoadPage()
	}
	
	doc, _ := gokogiri.ParseHtml([]byte(pc.Content))
	
	nodes, _ := doc.Search("//div[@style='padding-left:5px;padding-top:5px;']")
	
	reMatchId := regexp.MustCompile(`matchid=([0-9]+)`)
	reTeam := regexp.MustCompile(`([A-Za-z0-9\.\-_ ]+) \(([0-9]+)\)`)
	reTeamId := regexp.MustCompile(`teamid=([0-9]+)`)
	reEventId := regexp.MustCompile(`eventid=([0-9]+)`)
	for _, n := range nodes {
		
		
		// 0 -> match_id
		// 1 -> team 1
		// 2 -> team 2
		// 3 -> event
		links, _ := n.Search("./a")

		rs := reMatchId.FindStringSubmatch(links[0].Attribute("href").Value())

		_mId, _ := strconv.ParseInt(rs[1], 10, 32)
		dateStr := links[0].FirstChild().Content()
		
		rs = reTeam.FindStringSubmatch(strings.TrimSpace(links[1].FirstChild().Content()))
		team1Name := rs[1]
		_team1Score, _ := strconv.ParseInt(rs[2], 10, 32)
		
		rs = reTeamId.FindStringSubmatch(links[1].Attribute("href").Value())
		_team1Id, _ := strconv.ParseInt(rs[1], 10, 32)
		
		rs = reTeam.FindStringSubmatch(strings.TrimSpace(links[2].FirstChild().Content()))
		team2Name := rs[1]
		_team2Score, _ := strconv.ParseInt(rs[2], 10, 32)
		
		rs = reTeamId.FindStringSubmatch(links[2].Attribute("href").Value())
		_team2Id, _ := strconv.ParseInt(rs[1], 10, 32)
		
		rs = reEventId.FindStringSubmatch(links[3].Attribute("href").Value())
		_eventId, _ := strconv.ParseInt(rs[1], 10, 32)
		
		eventName := strings.TrimSpace(links[3].FirstChild().Content())
		
		match := &Match{}
		match.Date = ParseDate(dateStr)
		match.MatchId = int(_mId)
		match.Team1.Name = team1Name
		match.Team1.Score = int(_team1Score)
		match.Team1.TeamId = int(_team1Id)
		
		match.Team2.Name = team2Name
		match.Team2.Score = int(_team2Score)
		match.Team2.TeamId = int(_team2Id)
		
		match.Event.EventId = int(_eventId)
		match.Event.Name = eventName
		
		matches = append(matches, match)
	}
	
	doc.Free()
	
	return matches
	
}


func ParseDate(datestr string) MatchDate {
	
	re := regexp.MustCompile("([0-9]{1,2})/([0-9]{1,2}) ([0-9]{2})")
	rs := re.FindAllStringSubmatch(datestr, -1)
	
	day, _ := strconv.ParseInt(rs[0][1], 10, 32)
	month, _ := strconv.ParseInt(rs[0][2], 10, 32)
	year, _ := strconv.ParseInt(rs[0][3], 10, 32)
	
	year += 2000
	
	md := MatchDate{int(day), int(month), int(year)}
	return md
}

func StringStatToInt(stat string) int {
	int_stat := 0
	
	if stat != "-" { 
		_trans, _ := strconv.ParseInt(stat, 10, 32)
		int_stat = int(_trans)
	}
	
	return int_stat
}
