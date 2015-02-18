package csgoscrapper

import (
	"regexp"
	"strconv"
	"fmt"
)

type MatchTeam struct {
	TeamId int
	Score int
}

type MatchPlayerStat struct {
	PlayerId int
	TeamId int
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

type Match struct {
	MatchId int
	Date MatchDate
	Team1 MatchTeam
	Team2 MatchTeam
	Map string
	EventId int
	PlayerStats []*MatchPlayerStat
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
		
		t_id, _ := strconv.ParseInt(s[4], 10, 32)
		
		frags, _ := strconv.ParseInt(s[6], 10, 32)
		
		headshots := StringStatToInt(s[8])
		
		assists := StringStatToInt(s[10])
		
		deaths, _ := strconv.ParseInt(s[11], 10, 32)
		
		kdr, _ := strconv.ParseFloat(s[12], 32)
		
		kdrDelta, _ := strconv.ParseInt(s[14], 10, 32)
		
		rating, _ := strconv.ParseFloat(s[15], 32)
		
		stat := &MatchPlayerStat{int(p_id), int(t_id), int(frags), int(headshots), int(assists), int(deaths), float32(kdr), int(kdrDelta), float32(rating)}
		
		m.PlayerStats = append(m.PlayerStats, stat)
		
	}
	
	if len(m.PlayerStats) < 10 {
		log.Error(fmt.Sprintf("match[%d], only %d players, missings player !", m.MatchId, len(m.PlayerStats)))
		//output 
		//ioutil.WriteFile(os.TempDir() + "/gopool/" + strconv.Itoa(m.MatchId) + ".log", []byte(pc.Content), 0644)
	}
}