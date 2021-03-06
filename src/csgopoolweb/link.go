package csgopoolweb

import (
	"fmt"
	"strconv"
	"csgodb"
	"hltvscrapper"
	"eseascrapper"
)

type Parameter struct {
	Name string
	Value string
}

type Link struct {
	Caption string
	Url string
	Target string
	Params []Parameter
}

func GetTeamLink(t *csgodb.Team) string {
	html := ""
	
	if t.EseaId != 0 {
		href := eseascrapper.GetTeamURL(t.EseaId)
		html += fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, href, "ESEA.net")
	}
	
	if t.HltvId != 0  {
		if len(html) > 0 {
			html += " | "
		}
		
		pu := hltvscrapper.GetTeamPage(t.HltvId)
		href := pu.GenerateURL()
		html += fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, href, "HLTV.org")
	}
	
	return html
}

func GetPlayerLink(p *csgodb.Player) string {
	html := ""
	
	if p.EseaId != 0 {
		html += fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, eseascrapper.GetPlayerURL(p.EseaId), "ESEA.net")
	}
	
	if p.HltvId != 0 {
		p := hltvscrapper.GetPlayerPage(p.HltvId)
		href := p.GenerateURL()
		if len(html) > 0 {
			html += " | "
		}
		html += fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, href, "HLTV.org")
	}
	
	return html	
}

func GetMatchLink(m *csgodb.Match) string {
	html := ""
	href := ""
	
	if m.Source == csgodb.HltvSource {
		
		p := hltvscrapper.GetMatchPage(m.SourceId)
		href = p.GenerateURL()
		
	} else if m.Source == csgodb.EseaSource {
		href = eseascrapper.GetMatchURL(m.SourceId).Url()
	}
	
	html = fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, href, m.SourceName)
	
	return html
}

func GetMatchLinkCaption(m *csgodb.Match, caption string) string {
	html := ""
	href := ""
	
	if m.Source == csgodb.HltvSource {
		
		p := hltvscrapper.GetMatchPage(m.SourceId)
		href = p.GenerateURL()
		
	} else if m.Source == csgodb.EseaSource {
		href = eseascrapper.GetMatchURL(m.SourceId).Url()
	}
	
	html = fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, href, caption)
	
	return html
}

func (l *Link) AddParameter(name string, value string) {
	
	p := Parameter{Name: name, Value: value}
	l.Params = append(l.Params, p)
}

func (l *Link) AddInt(name string, value int) {
   
  l.AddParameter(name, strconv.Itoa(value))
}

func (l *Link) GetHTML() string {
	full_url := l.Url
	
	if len(l.Params) > 0 {
		full_url = full_url + "?"
		
		for i, p := range l.Params {
			if i == 0 {
				full_url = full_url + fmt.Sprintf("%s=%s", p.Name, p.Value)
			} else {
				full_url = full_url + fmt.Sprintf("&%s=%s", p.Name, p.Value)
			}
		}
		
	}
	
	target := ""
	
	if len(l.Target) > 0 {
		target = fmt.Sprintf(` target="%s"`, l.Target)
	}
	
	return fmt.Sprintf(`<a href="%s"%s>%s</a>`, full_url, target, l.Caption)
}