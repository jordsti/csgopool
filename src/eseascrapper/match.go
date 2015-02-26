package eseascrapper

import (
	"regexp"
	"strconv"
	"fmt"
	"github.com/moovweb/gokogiri"
	"strings"
)

const (
	Completed = 1
	Forfeit = 2
)

type MatchDate struct {
	Year int
	Month int
	Day int
}

type MatchTeam struct {
	TeamId int
	Score int
	Name string
}

type PlayerMatchStat struct {
	PlayerId int
	TeamId int
	Name string
	RWS float32
	Frags int
	Assists int
	Deaths int
	BombPlants int
	BombDefusal int
	RoundPlayed int
	
	KDRatio float32
	KDDelta int
}

type Match struct {
	MatchId int
	Date MatchDate
	Map string
	Team1 *MatchTeam
	Team2 *MatchTeam
	PlayerStats []*PlayerMatchStat
	Status int
}

func (pms *PlayerMatchStat) Player() Player {
	return Player{PlayerId: pms.PlayerId, Name: pms.Name}
}

func (mt *MatchTeam) Team() *Team {
	return &Team{Name: mt.Name, TeamId: mt.TeamId}
}

func (pc *PageContent) ParseMatches() []*Match {
	matches := []*Match{}
	
	doc, _ := gokogiri.ParseHtml([]byte(pc.Content))
	
	nodes, _ := doc.Search("//div[@class='match-container']")
	
	re := regexp.MustCompile(`^/index.php\?s=stats&d=match&id=([0-9]+)$`)
	for _, n := range nodes {
		
		nOverview, _ := n.Search("./div[@class='match-overview']")
		footer, _ := n.Search("./div[@class='match-footer']")

		innerNodes, _ := nOverview[0].Search("./table/tr/th")
		status := 0
		if len(innerNodes) > 0 {
			matchStatus := strings.TrimSpace(innerNodes[0].Content())
			fmt.Printf("Match Status : %s\n", matchStatus)
			if matchStatus == "Completed" {
				status = Completed
				
			} else if matchStatus == "Completed (Forfeit)" {
				status = Forfeit
			}
		}

		if status != 0 {
			linkNode, _ := nOverview[0].Search("./a")
			link := linkNode[0].Attribute("href").Value()
			rs := re.FindStringSubmatch(link)
			_m_id, _ := strconv.ParseInt(rs[1], 10, 32)
			matchId := int(_m_id)
			fmt.Printf("Match Id : %d\n", matchId)
			
			data := strings.Split(footer[0].Content(), "/")
			
			mapstr := strings.TrimSpace(data[1])
			
			m := &Match{}
			m.Status = status
			m.MatchId = matchId
			m.Date.Year = pc.Url.Date.Year()
			m.Date.Month = int(pc.Url.Date.Month())
			m.Date.Day = pc.Url.Date.Day()
			m.Map = mapstr
			matches = append(matches, m)
		}
		
	}
	
	doc.Free()
	return matches
}

func (m *Match) ParseMatch() {
	
	url := GetMatchURL(m.MatchId)
	pc := url.LoadPage()
	
	re_teamid := regexp.MustCompile("^/teams/([0-9]+)$")
	//fmt.Println(pc.Content)
	
	//parsing team name
	doc, _ := gokogiri.ParseHtml([]byte(pc.Content))
	nodes, _ := doc.Search("//div[@id='body-match-stats']/table[@class='box']")
	scoresRow, _ := nodes[0].Search("./tr")
	
	//row -> 1, Team1
	//row -> 2, Team2
	
	if len(scoresRow) >= 3 {
		nodes, _ = scoresRow[1].Search("./th[@align='left']/a")
		team1 := nodes[0]
		rs := re_teamid.FindStringSubmatch(team1.Attribute("href").Value())
		
		_team_id, _ := strconv.ParseInt(rs[1], 10, 32)
		//fmt.Printf("%v\n", rs[1])
		m.Team1 = &MatchTeam{TeamId: int(_team_id), Name: team1.Content()}
		
		//scores parsing
		
		nodes, _ = scoresRow[1].Search("./td[@class='ct stat']")
		_score, _ := strconv.ParseInt(nodes[0].Content(), 10, 32)
		m.Team1.Score = int(_score)
		
		nodes, _ = scoresRow[1].Search("./td[@class='t stat']")
		_score, _ = strconv.ParseInt(nodes[0].Content(), 10, 32)
		m.Team1.Score += int(_score)
		
		nodes, _ = scoresRow[2].Search("./th[@align='left']/a")
		team2 := nodes[0]
		rs = re_teamid.FindStringSubmatch(team2.Attribute("href").Value())
		
		_team_id, _ = strconv.ParseInt(rs[1], 10, 32)
		m.Team2 = &MatchTeam{TeamId: int(_team_id), Name: team2.Content()}

		
		//scores parsing
		
		nodes, _ = scoresRow[2].Search("./td[@class='ct stat']")
		_score, _ = strconv.ParseInt(nodes[0].Content(), 10, 32)
		m.Team2.Score = int(_score)
		
		nodes, _ = scoresRow[2].Search("./td[@class='t stat']")
		_score, _ = strconv.ParseInt(nodes[0].Content(), 10, 32)
		m.Team2.Score += int(_score)
	}
	
	re := regexp.MustCompile(`^/users/([0-9]+)$`)
	
	//team1
	els, _ := doc.Search("//tbody[@id='body-match-total1']/tr")

	for _, el := range els {
		links, _  := el.Search("./td/a")
		
		plink := links[1]
		rs := re.FindAllStringSubmatch(plink.Attribute("href").Value(), -1)	
		playerId, _ := strconv.ParseInt(rs[0][1], 10, 32)
		fmt.Printf("ESEA Player Id :%d\n", playerId)
		
		stats, _ := el.Search("./td[@class='stat']")
		rws, _ := strconv.ParseFloat(stats[0].Content(), 32)
		frags, _ := strconv.ParseInt(stats[1].Content(), 10, 32)
		assists, _ := strconv.ParseInt(stats[2].Content(), 10, 32)
		deaths, _ := strconv.ParseInt(stats[3].Content(), 10, 32)
		bombPlants, _ := strconv.ParseInt(stats[4].Content(), 10, 32)
		bombDefusal, _ := strconv.ParseInt(stats[5].Content(), 10, 32)
		roundPlayed, _ := strconv.ParseInt(stats[6].Content(), 10, 32)

		pstat := &PlayerMatchStat{}
		pstat.PlayerId = int(playerId)
		pstat.TeamId = m.Team1.TeamId
		pstat.Name = plink.Content()
		pstat.RWS = float32(rws)
		pstat.Frags = int(frags)
		pstat.Assists = int(assists)
		pstat.Deaths = int(deaths)
		pstat.BombPlants = int(bombPlants)
		pstat.BombDefusal = int(bombDefusal)
		pstat.RoundPlayed = int(roundPlayed)
		
		pstat.KDRatio = float32(pstat.Frags) / float32(pstat.Deaths)
		pstat.KDDelta = pstat.Frags - pstat.Deaths
		
		m.PlayerStats = append(m.PlayerStats, pstat)
	}
	
	//team2
	els, _ = doc.Search("//tbody[@id='body-match-total2']/tr")

	for _, el := range els {
		links, _  := el.Search("./td/a")
		
		plink := links[1]
		rs := re.FindAllStringSubmatch(plink.Attribute("href").Value(), -1)
		
		playerId, _ := strconv.ParseInt(rs[0][1], 10, 32)
		
		stats, _ := el.Search("./td[@class='stat']")
		rws, _ := strconv.ParseFloat(stats[0].Content(), 32)
		frags, _ := strconv.ParseInt(stats[1].Content(), 10, 32)
		assists, _ := strconv.ParseInt(stats[2].Content(), 10, 32)
		deaths, _ := strconv.ParseInt(stats[3].Content(), 10, 32)
		bombPlants, _ := strconv.ParseInt(stats[4].Content(), 10, 32)
		bombDefusal, _ := strconv.ParseInt(stats[5].Content(), 10, 32)
		roundPlayed, _ := strconv.ParseInt(stats[6].Content(), 10, 32)

		pstat := &PlayerMatchStat{}
		pstat.PlayerId = int(playerId)
		pstat.TeamId = m.Team2.TeamId
		pstat.Name = plink.Content()
		pstat.RWS = float32(rws)
		pstat.Frags = int(frags)
		pstat.Assists = int(assists)
		pstat.Deaths = int(deaths)
		pstat.BombPlants = int(bombPlants)
		pstat.BombDefusal = int(bombDefusal)
		pstat.RoundPlayed = int(roundPlayed)
		
		pstat.KDRatio = float32(pstat.Frags) / float32(pstat.Deaths)
		pstat.KDDelta = pstat.Frags - pstat.Deaths
		
		m.PlayerStats = append(m.PlayerStats, pstat)
	}
	doc.Free()
}
