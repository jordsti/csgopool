package csgoscrapper

import (
	"regexp"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type Event struct {
	EventId int
	Name string
	Matches []Match
}

func GetMatchById(events []*Event, id int) *Match {
	
	for _, evt := range events {
		for _, m := range evt.Matches {
			if m.MatchId == id {
				return &m
			}
		}
	}
	
	return nil
}

func GetEventById(events []*Event, id int) *Event {
	
	for _, evt := range events {
		if evt.EventId == id {
			return evt
		}
	}
	
	return nil
}

func EventExists(events []*Event, id int) bool {
	
	for _, evt := range events {
		if evt.EventId == id {
			return true
		}
	}
	
	return false
}

func GetLastEvent(events []*Event) *Event {
	lastId := 0
	curEvent := &Event{}
	for _, evt := range events {
		if evt.EventId > lastId {
			curEvent = evt
			lastId = evt.EventId
		}
	}
	
	return curEvent
}

func (p PageContent) ParseEventsWithoutMatches() []*Event {
	events := []*Event{}
	
	re := regexp.MustCompile("<a href=\"/\\?pageid=183&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:50%;float:left;\"><img style=\"vertical-align: -1px;\" src=\"http://static.hltv.org//images/mod_csgo.png\" title=\"Counter-Strike: Global Offensive\"> (.+)</div></a>")
	rs := re.FindAllStringSubmatch(p.Content, -1)
	
	for _, evt := range rs {
		
		e_id, _ := strconv.ParseInt(evt[1], 10, 32)
		e_name := evt[2]
		
		e := &Event{EventId: int(e_id), Name: e_name}
		
		//e.LoadAllMatches()
		

		events = append(events, e)

	}
	
	return events
}

func (p PageContent) ParseEvents() []*Event {
	events := []*Event{}
	
	re := regexp.MustCompile("<a href=\"/\\?pageid=183&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:50%;float:left;\"><img style=\"vertical-align: -1px;\" src=\"http://static.hltv.org//images/mod_csgo.png\" title=\"Counter-Strike: Global Offensive\"> (.+)</div></a>")
	rs := re.FindAllStringSubmatch(p.Content, -1)
	
	for _, evt := range rs {
		
		e_id, _ := strconv.ParseInt(evt[1], 10, 32)
		e_name := evt[2]
		
		e := &Event{EventId: int(e_id), Name: e_name}
		
		e.LoadAllMatches()
		
		if len(e.Matches) > 0 {
			events = append(events, e)
		} else {
			log.Info("Breaking events parsing")
			break
		}
	}
	
	return events
}

func (p PageContent) ParseMatches() []*Match {
	//parsing all matches
	//match id, event id, date, team id 1, event id, flag1, team name 1, score 1, team id 2, event id, flag 2, team name 2, score 2, map
	matches := []*Match{}
	
	re := regexp.MustCompile("<a href=\"/\\?pageid=188&amp;matchid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:10%;float:left;;font-weight:normal;\">([0-9/ ]+)</div></a><a href=\"/\\?pageid=179&amp;teamid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:25%;float:left;;font-weight:normal;\"><img style=\"vertical-align:-20%;\" src=\"(.+)\" alt=\"\" height=\"12\" width=\"18\" class=\"flagFix\"/> (.+) \\(([0-9]+)\\)</div></a><a href=\"/\\?pageid=179&amp;teamid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:25%;float:left;;font-weight:normal;\"><img style=\"vertical-align:-20%;\" src=\"(.+)\" alt=\"\" height=\"12\" width=\"18\" class=\"flagFix\"/> (.+) \\(([0-9]+)\\)</div></a><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;text-align:center;font-weight:normal;color:black;\">([a-z0-9]+)</div><a href=\"/\\?pageid=188&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:30%;float:left;font-weight:normal;\"><img style=\"vertical-align: -1px;\" src=\"http://static.hltv.org//images/mod_csgo.png\" title=\"Counter-Strike: Global Offensive\"> <span title=\"(.+)\">(.+)</span></div></a>")
	rs := re.FindAllStringSubmatch(p.Content, -1)
	
	for _, m := range rs {
		
		m_id, _ := strconv.ParseInt(m[1], 10, 32)

		log.Info(fmt.Sprintf("Parsing match [%d]", m_id))

		date := m[3]
		t_1, _ := strconv.ParseInt(m[4], 10, 32)
		ts_1, _ := strconv.ParseInt(m[8], 10, 32)
		
		t_2, _ := strconv.ParseInt(m[9], 10, 32)
		ts_2, _ := strconv.ParseInt(m[13], 10, 32)
		
		gameMap := m[14]
		
		match := &Match{MatchId: int(m_id), Date: ParseDate(date), Map: gameMap}
		
		match.Team1.TeamId = int(t_1)
		match.Team1.Score = int(ts_1)
		
		match.Team2.TeamId = int(t_2)
		match.Team2.Score = int(ts_2)
		
		
		matches = append(matches, match)
	}
	
	return matches
}

func (p PageContent) UpdateEvents(events []*Event) []*Event {

	re := regexp.MustCompile("<a href=\"/\\?pageid=183&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:50%;float:left;\"><img style=\"vertical-align: -1px;\" src=\"http://static.hltv.org//images/mod_csgo.png\" title=\"Counter-Strike: Global Offensive\"> (.+)</div></a>")
	rs := re.FindAllStringSubmatch(p.Content, -1)
	
	for _, evt := range rs {
		
		e_id, _ := strconv.ParseInt(evt[1], 10, 32)
		e_name := evt[2]
		
		if EventExists(events, int(e_id)) {
			//updating the matches
			e := GetEventById(events, int(e_id))
			if e != nil {
				log.Info(fmt.Sprintf("Update event [%d]", e.EventId))
				e.UpdateMatches()
			} else {
				log.Error(fmt.Sprintf("Fatal error event [%d] not found", int(e_id)))
			}
			
		} else {
			log.Info(fmt.Sprintf("Adding a new event [%d]", int(e_id)))
			
			e := &Event{EventId: int(e_id), Name: e_name}
		
			e.LoadAllMatches()
			
			if len(e.Matches) > 0 {
				events = append(events, e)
			} else {
				log.Info("Breaking events parsing")
				break
			}
		}
		
	}
	
	return events
}

func (e *Event) UpdateMatches() {
	
	page := GetEventMatches(e.EventId)
	pc, _ := page.LoadPage()
	
	log.Info(fmt.Sprintf("Update Event [%d], Status [%d]", e.EventId, pc.Status))
	
	//parsing all matches
	//match id, event id, date, team id 1, event id, flag1, team name 1, score 1, team id 2, event id, flag 2, team name 2, score 2, map
	re := regexp.MustCompile("<a href=\"/\\?pageid=188&amp;matchid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:10%;float:left;;font-weight:normal;\">([0-9/ ]+)</div></a><a href=\"/\\?pageid=179&amp;teamid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:25%;float:left;;font-weight:normal;\"><img style=\"vertical-align:-20%;\" src=\"(.+)\" alt=\"\" height=\"12\" width=\"18\" class=\"flagFix\"/> (.+) \\(([0-9]+)\\)</div></a><a href=\"/\\?pageid=179&amp;teamid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:25%;float:left;;font-weight:normal;\"><img style=\"vertical-align:-20%;\" src=\"(.+)\" alt=\"\" height=\"12\" width=\"18\" class=\"flagFix\"/> (.+) \\(([0-9]+)\\)</div></a><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;text-align:center;font-weight:normal;color:black;\">([a-z0-9]+)</div><a href=\"/\\?pageid=188&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:30%;float:left;font-weight:normal;\"><img style=\"vertical-align: -1px;\" src=\"http://static.hltv.org//images/mod_csgo.png\" title=\"Counter-Strike: Global Offensive\"> <span title=\"(.+)\">(.+)</span></div></a>")
	rs := re.FindAllStringSubmatch(pc.Content, -1)
	
	for _, m := range rs {
		
		m_id, _ := strconv.ParseInt(m[1], 10, 32)
		if e.MatchExists(int(m_id)) {
			
			log.Info(fmt.Sprintf("Match [%d] already exists in the event [%d]", m_id, e.EventId))
			
		} else {
			
			log.Info(fmt.Sprintf("Adding match [%d] to the event [%d]", m_id, e.EventId))
			date := m[3]
			t_1, _ := strconv.ParseInt(m[4], 10, 32)
			ts_1, _ := strconv.ParseInt(m[8], 10, 32)
			
			t_2, _ := strconv.ParseInt(m[9], 10, 32)
			ts_2, _ := strconv.ParseInt(m[13], 10, 32)
			
			gameMap := m[14]
			
			match := Match{MatchId: int(m_id), Date: ParseDate(date), Map: gameMap}
			//todo global to edit this
			if match.Date.Year < State.MinYear {
				//skipping this match, too old
				log.Info(fmt.Sprintf("Skipping match[%d], was too old, year %d\n", match.MatchId, match.Date.Year))
			} else {
			
				match.Team1.TeamId = int(t_1)
				match.Team1.Score = int(ts_1)
				
				match.Team2.TeamId = int(t_2)
				match.Team2.Score = int(ts_2)
				
				m_ptr := &match
				m_ptr.GetMatchStats()
				
				e.Matches = append(e.Matches, match)
			}	
		}
	}
}

func (e *Event) MatchExists(matchId int) bool {
	
	for _, m := range e.Matches {
		if m.MatchId == matchId {
			return true
		}
	}
	
	return false
	
}

func (e *Event) LoadAllMatches() {
	
	page := GetEventMatches(e.EventId)
	pc, _ := page.LoadPage()
	log.Info(fmt.Sprintf("Fetching Event [%d], Status [%d]", e.EventId, pc.Status))
	//parsing all matches
	//match id, event id, date, team id 1, event id, flag1, team name 1, score 1, team id 2, event id, flag 2, team name 2, score 2, map
	re := regexp.MustCompile("<a href=\"/\\?pageid=188&amp;matchid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:10%;float:left;;font-weight:normal;\">([0-9/ ]+)</div></a><a href=\"/\\?pageid=179&amp;teamid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:25%;float:left;;font-weight:normal;\"><img style=\"vertical-align:-20%;\" src=\"(.+)\" alt=\"\" height=\"12\" width=\"18\" class=\"flagFix\"/> (.+) \\(([0-9]+)\\)</div></a><a href=\"/\\?pageid=179&amp;teamid=([0-9]+)&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:25%;float:left;;font-weight:normal;\"><img style=\"vertical-align:-20%;\" src=\"(.+)\" alt=\"\" height=\"12\" width=\"18\" class=\"flagFix\"/> (.+) \\(([0-9]+)\\)</div></a><div class=\"covSmallHeadline\" style=\"font-weight:normal;width:10%;float:left;text-align:center;font-weight:normal;color:black;\">([a-z0-9]+)</div><a href=\"/\\?pageid=188&amp;eventid=([0-9]+)&amp;gameid=2\"><div class=\"covSmallHeadline\" style=\"width:30%;float:left;font-weight:normal;\"><img style=\"vertical-align: -1px;\" src=\"http://static.hltv.org//images/mod_csgo.png\" title=\"Counter-Strike: Global Offensive\"> <span title=\"(.+)\">(.+)</span></div></a>")
	rs := re.FindAllStringSubmatch(pc.Content, -1)
	
	for _, m := range rs {
		
		m_id, _ := strconv.ParseInt(m[1], 10, 32)
		date := m[3]
		t_1, _ := strconv.ParseInt(m[4], 10, 32)
		ts_1, _ := strconv.ParseInt(m[8], 10, 32)
		
		t_2, _ := strconv.ParseInt(m[9], 10, 32)
		ts_2, _ := strconv.ParseInt(m[13], 10, 32)
		
		gameMap := m[14]
		
		match := Match{MatchId: int(m_id), Date: ParseDate(date), Map: gameMap}
		//todo global to edit this
		if match.Date.Year < State.MinYear {
			//skipping this match, too old
			log.Info(fmt.Sprintf("Skipping match[%d], was too old, year %d", match.MatchId, match.Date.Year))
		} else {
		
		match.Team1.TeamId = int(t_1)
		match.Team1.Score = int(ts_1)
		
		match.Team2.TeamId = int(t_2)
		match.Team2.Score = int(ts_2)
		
		m_ptr := &match
		m_ptr.GetMatchStats()
		
		e.Matches = append(e.Matches, match)
		
		}
	}
}

func SaveEvents(events []*Event, path string) {
	b, err := json.MarshalIndent(events, "", "	")
	if err != nil {
		log.Error("Error while saving events [1]")
	}
	
	err = ioutil.WriteFile(path, b, 0644)
	
	if err != nil {
		log.Error("Error while saving events [2]")
	}
}

func LoadEvents(path string) []*Event {
	events := []*Event{}
	
	data, err := ioutil.ReadFile(path)
	
	if err != nil {
		log.Error("Error while reading events [1]")
	}
	
	err = json.Unmarshal(data, &events)
	
	if err != nil {
		log.Error("Error while reading events [2]")
		//fmt.Println(err)
	}
	
	return events
}
