package eseascrapper

import (
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
)

var eseaUrl = "http://play.esea.net/"

const (
	NorthAmerica = 1
	Europe = 2
	Brazil = 4
	AsiaPacific = 5
	Australia = 6
	India = 7
	CSGOId = 25
	)

type PageURL struct {
	RegionId int
	GameId int
	Id int
	Date time.Time
	Division string
	Section string
	Page string
	Action string
}

type PageContent struct {
	Url *PageURL
	Content string
	Status int
}

func GetScheduleURL(date time.Time, division string, region int) *PageURL{
	
	url := &PageURL{}
	
	url.Date = date
	url.RegionId = region
	url.GameId = CSGOId
	url.Division = division
	url.Section = "league"
	url.Page = "index.php"
	url.Action = "schedule"
	
	return url
}

func GetMatchURL(matchId int) *PageURL {
	
	url := &PageURL{}
	
	url.Id = matchId
	url.Section = "stats"
	url.Action = "match"
	
	return url
}

func (p *PageURL) Url() string {
	
	date_str := ""
	
	if p.Date.Year() > 2000 {
		date_str = fmt.Sprintf("%d-%02d-%02d", p.Date.Year(), p.Date.Month(), p.Date.Day())
	}
	
	url := eseaUrl
	url += p.Page + "?"
	
	if len(p.Section) > 0 {
		url += "s=" + p.Section
	}
	
	if len(p.Action) > 0 {
		url += "&d=" + p.Action
	}
	
	if p.GameId != 0 {
		url += fmt.Sprintf("&game_id=%d", p.GameId)
	}
	
	if len(p.Division) > 0 {
		url += "&division_level=" + p.Division
	}
	
	if p.RegionId != 0 {
		url += fmt.Sprintf("&region_id=%d", p.RegionId)
	}
	
	if len(date_str) > 0 {
		url += "&date=" + date_str
	}
	
	if p.Id != 0 {
		url += fmt.Sprintf("&id=%d", p.Id)
	}
	
	return url
}

func (p *PageURL) LoadPage() *PageContent {
	
	content := &PageContent{Url: p}
	
	resp, err := http.Get(p.Url())
	
	defer resp.Body.Close()
	
	if err != nil {
		fmt.Printf("Error while loading %s", p.Url())
	}
	
	body, err2 := ioutil.ReadAll(resp.Body)
	
	if err2 != nil {
		fmt.Printf("Error while loading %s", p.Url())
	}
	
	content.Content = string(body)
	content.Status = resp.StatusCode
	
	return content
}