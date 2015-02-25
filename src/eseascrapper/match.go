package eseascrapper

import (
	"regexp"
	"strconv"

	"github.com/moovweb/gokogiri"
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
}

type Match struct {
	MatchId int
	Date MatchDate
	Team1 MatchTeam
	Team2 MatchTeam
	PlayerStats []*PlayerMatchStat
}


func (pc *PageContent) ParseMatches() []*Match {
	matches := []*Match{}
	
	re := regexp.MustCompile(`<a href="/index.php\?s=stats&d=match&id=([0-9]+)">Statistics and Discussion</a>`)
	rs := re.FindAllStringSubmatch(pc.Content, -1)
	
	for _, mrs := range rs {

		m := &Match{}
		_match_id, _ := strconv.ParseInt(mrs[1], 10, 32)
		
		m.MatchId = int(_match_id)
		m.Date.Year = pc.Url.Date.Year()
		m.Date.Month = int(pc.Url.Date.Month())
		m.Date.Day = pc.Url.Date.Day()
		
		matches = append(matches, m)
	}
	
	return matches
}

func (m *Match) ParseMatch() {
	
	url := GetMatchURL(m.MatchId)
	pc := url.LoadPage()
	
	re := regexp.MustCompile(`<th align="left"><a href="/teams/([0-9]+)">([A-Za-z0-9\-_ \.]+)</a></th>\s+<td class="(ct|t) stat">([0-9]+)</td>\s+<td class="(ct|t) stat">([0-9]+)</td>`)
	rs := re.FindAllStringSubmatch(pc.Content, -1)
	iteam := 0
	for _, mrs := range rs {
		_team_id, _ := strconv.ParseInt(mrs[1], 10, 32)
		_team_name := mrs[2]
		
		_stat_1, _ := strconv.ParseInt(mrs[4], 10, 32)
		_stat_2, _ := strconv.ParseInt(mrs[6], 10, 32)
		
		if iteam == 0 {
			m.Team1.TeamId = int(_team_id)
			m.Team1.Name = _team_name
			m.Team1.Score = int(_stat_1 + _stat_2)
		} else if iteam == 1 {
			m.Team2.TeamId = int(_team_id)
			m.Team2.Name = _team_name
			m.Team2.Score = int(_stat_1 + _stat_2)
		}
		
		iteam++
	}
	
	re = regexp.MustCompile(`^/users/([0-9]+)$`)
	doc, _ := gokogiri.ParseHtml([]byte(pc.Content))
	
	//team1
	els, _ := doc.Search("//tbody[@id='body-match-total1']/tr")

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
		pstat.TeamId = m.Team1.TeamId
		pstat.Name = plink.Content()
		pstat.RWS = float32(rws)
		pstat.Frags = int(frags)
		pstat.Assists = int(assists)
		pstat.Deaths = int(deaths)
		pstat.BombPlants = int(bombPlants)
		pstat.BombDefusal = int(bombDefusal)
		pstat.RoundPlayed = int(roundPlayed)
		
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
		
		m.PlayerStats = append(m.PlayerStats, pstat)
	}
	doc.Free()
}
