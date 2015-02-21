package csgopoolweb

import (
	"fmt"
	"strconv"
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