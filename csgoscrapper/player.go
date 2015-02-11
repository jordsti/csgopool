package csgoscrapper

import (
	"regexp"
	"strconv"
	"fmt"
)

type PlayerStat struct {
	Frags int
	Headshot float32
	Deaths int
	KDRatio float32
	MapsPlayed int
	RoundsPlayed int
	AvgFragsPerRound float32
	AvgAssistsPerRound float32
	AvgDeathsPerRound float32
	Rating float32
}

type Player struct {
	Name string
	PlayerId int
	Stats PlayerStat
}

func (p PageContent) ParsePlayer() []Player {
	
	players := []Player{}
	
	re := regexp.MustCompile("<a href=\"/\\?pageid=173&amp;playerid=([0-9]+)&amp;eventid=0&amp;gameid=2\">(.+)</a>")
	rs := re.FindAllStringSubmatch(p.Content, -1)
	
	for _, p := range rs {
		//fmt.Printf("%d\n", p[1])
		p_id, _ := strconv.ParseInt(p[1], 10, 32)
		player := Player{Name: p[2], PlayerId: int(p_id)}
		
		pl := &player
		pl.LoadStats()
		
		players = append(players, player)
	}
	
	return players
}

func (p *Player) LoadStats() {
	
	//generating url
	page := GetPlayerPage(p.PlayerId)
	content, _ := page.LoadPage()
	
	if len(p.Name) == 0 {
		//parsing name too
		fmt.Printf("Parsing name for player [%d]\n", p.PlayerId)
		ren := regexp.MustCompile("Player stats: ([a-zA-Z0-9\\.\\-_ ]+) <span class=\"tab_spacer\">")
		rsn := ren.FindAllStringSubmatch(content.Content, -1)
		p.Name = rsn[0][1]
	}
	
	//total kills
	re := regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Total kills</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+)</div>")
	rs := re.FindAllStringSubmatch(content.Content, -1)
	
	kills, _ := strconv.ParseInt(rs[0][1], 10, 32)
	
	p.Stats.Frags = int(kills)
	
	//Headshot %
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Headshot %</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+\\.?[0-9]*)%</div>")
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	headshots, _ := strconv.ParseFloat(rs[0][1], 32)
	
	p.Stats.Headshot = float32(headshots)
	
	//Total deaths
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Total deaths</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+)</div>")
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	deaths, _ := strconv.ParseInt(rs[0][1], 10, 32)
	
	p.Stats.Deaths = int(deaths)
	
	//K/D ratio
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">K/D Ratio</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+\\.?[0-9]*)</div>")
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	kdr, _ := strconv.ParseFloat(rs[0][1], 32)
	
	p.Stats.KDRatio = float32(kdr)
	
	//Maps Played
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Maps played</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+)</div>")	
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	mapsPlayed, _ := strconv.ParseInt(rs[0][1], 10, 32)
	
	p.Stats.MapsPlayed = int(mapsPlayed)
	
	//Rounds played
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Rounds played</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+)</div>")
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	rounds, _ := strconv.ParseInt(rs[0][1], 10, 32)
	
	p.Stats.RoundsPlayed = int(rounds)
	
	//Avg Kills per round
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Average kills per round</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+\\.?[0-9]*)</div>")
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	avgkills, _ := strconv.ParseFloat(rs[0][1], 32)
	
	p.Stats.AvgFragsPerRound = float32(avgkills)
	
	//Avg assists per round
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Average assists per round</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+\\.?[0-9]*)</div>")
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	avgassists, _ := strconv.ParseFloat(rs[0][1], 32)
	
	p.Stats.AvgAssistsPerRound = float32(avgassists)
	
	//Avg deaths per round
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Average deaths per round</div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black\">([0-9]+\\.?[0-9]*)</div>")
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	avgdeaths, _ := strconv.ParseFloat(rs[0][1], 32)
	
	p.Stats.AvgDeathsPerRound = float32(avgdeaths)
	
	//Rating
	re = regexp.MustCompile("<div class=\"covSmallHeadline\" style=\"font-weight:normal;width:185px;float:left;text-align:left;font-weight:bold\">Rating <a href=\"/\\?pageid=242\" style=\"color:black;font-weight:normal\" title=\"Click here to see how rating is calculated\">\\(\\?\\)</a></div><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:100px;float:left;text-align:right;color:black;font-weight:bold\">([0-9]+\\.?[0-9]*)</div>")
	rs = re.FindAllStringSubmatch(content.Content, -1)
	
	rating, _ := strconv.ParseFloat(rs[0][1], 32)
	
	p.Stats.Rating = float32(rating)
}