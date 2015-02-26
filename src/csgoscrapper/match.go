package csgoscrapper

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
	
	fmt.Printf("Parse match [%d] from hltv.org\n", m.MatchId)
	
	url := GetMatchPage(m.MatchId)
	pc, _ := url.LoadPage()
	//error handling
	
	doc, _ := gokogiri.ParseHtml([]byte(pc.Content))
	
	//parse map
	nodes, _ := doc.Search("//div[@class='covSmallHeadline' and @style='font-weight:normal;width:180px;float:left;text-align:right;']")
	//first node is the map
	m.Map = strings.TrimSpace(nodes[0].Content())
	
	//player and stats
	
	nodes, _ = doc.Search("//div[starts-with(@style,'width:606px;height:22px;background-color:')]/div[@style='padding-left:5px;padding-top:5px;']")
	
	rePlayerId := regexp.MustCompile(`/\?pageid=173&playerid=([0-9]+)&gameid=2`)
	reTeamId := regexp.MustCompile(`/\?pageid=179&teamid=([0-9]+)`)
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
		ms.KDRatio = float32(frags) / float32(deaths)
		
		m.PlayerStats = append(m.PlayerStats, ms)
	}	
}

func GetMatches(offset int) []*Match {
	matches := []*Match{}
	
	url := GetMatchesPage(offset)
	//error handling
	pc, _ := url.LoadPage()
	//fmt.Printf(pc.Content)
	doc, _ := gokogiri.ParseHtml([]byte(pc.Content))
	
	nodes, _ := doc.Search("//div[@style='padding-left:5px;padding-top:5px;']")
	
	reMatchId := regexp.MustCompile(`/\?pageid=188&matchid=([0-9]+)&eventid=0&gameid=2`)
	reTeam := regexp.MustCompile(`([A-Za-z0-9\.\-_ ]+) \(([0-9]+)\)`)
	reTeamId := regexp.MustCompile(`/\?pageid=179&teamid=([0-9]+)&eventid=0&gameid=2`)
	reEventId := regexp.MustCompile(`/\?pageid=188&eventid=([0-9]+)&gameid=2`)
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

func (m *Match) IsPlayerIn(playerId int) bool {
  for _, pl := range m.PlayerStats {
    if pl.PlayerId == playerId {
      return true
    }
    
  }
  return false
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

func (m *Match) GetMatchStats() {
	page := GetMatchPage(m.MatchId)
	pc, _ := page.LoadPage()
	
	for pc.Status != 200 {
		log.Error(fmt.Sprintf("Match [%d], Status [%d], new attempt", m.MatchId, pc.Status))
		pc, _ = page.LoadPage()
	}
	
	log.Info(fmt.Sprintf("Match [%d], Status [%d]", m.MatchId, pc.Status))
	// 1 -> Flag
	// 2 -> player id
	// 3 -> player name
	// 4 -> team id
	// 5 -> team name
	// 6 -> frags
	// 7 
	// 8 -> headshot
	// 9
	// 10 -> assists
	// 11 -> deaths
	// 12 -> k/d
	// 13 -> kd color
	// 14 -> k/d delta
	// 15 -> rating

	re := regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:20%;float:left;\"><img src=\"(.+)\" alt=\"\" /> <a href=\"/\\?pageid=173&amp;playerid=([0-9]+)&amp;gameid=2\">(.+)</a></div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:20%;float:left;text-align:center\"><a href=\"/\\?pageid=179&amp;teamid=([0-9]+)\">(.+)</a></div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;text-align:center\">([0-9]+) (<span title=\"headshots\" style=\"cursor:help\">)?\\((.+?)\\)(</span>)?</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;text-align:center\">(.+?)</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;text-align:center\">([0-9]+)</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;text-align:center\">([0-9\\-]+\\.?[0-9]*)</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;color:(.+);text-align:center\">(.+)</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;text-align:center\">([0-9]\\.?[0-9]*)</div>")
	rs := re.FindAllStringSubmatch(pc.Content, -1)
	

	for _, s := range rs {
		
		p_id, _ := strconv.ParseInt(s[2], 10, 32)
		p_name := s[3]
		
		t_id, _ := strconv.ParseInt(s[4], 10, 32)
		
		frags, _ := strconv.ParseInt(s[6], 10, 32)
		
		headshots := StringStatToInt(s[8])
		
		assists := StringStatToInt(s[10])
		
		deaths, _ := strconv.ParseInt(s[11], 10, 32)
		
		kdr, _ := strconv.ParseFloat(s[12], 32)
		
		kdrDelta, _ := strconv.ParseInt(s[14], 10, 32)
		
		rating, _ := strconv.ParseFloat(s[15], 32)
		
		//stat := &MatchPlayerStat{int(p_id), p_name, int(t_id), int(frags), int(headshots), int(assists), int(deaths), float32(kdr), int(kdrDelta), float32(rating)}
		
		stat := &MatchPlayerStat{int(p_id), p_name, int(t_id), "", int(frags), int(headshots), int(assists), int(deaths), float32(kdr), int(kdrDelta), float32(rating)}
		
		m.PlayerStats = append(m.PlayerStats, stat)
		
	}
	
	if len(m.PlayerStats) < 10 {
		log.Error(fmt.Sprintf("match[%d], only %d players, missings player !", m.MatchId, len(m.PlayerStats)))
		//output 
		//ioutil.WriteFile(os.TempDir() + "/gopool/" + strconv.Itoa(m.MatchId) + ".log", []byte(pc.Content), 0644)
	}
}