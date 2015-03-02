package steamapi

import (
	"net/http"
	"fmt"
	"net/url"
	"bytes"
	"io/ioutil"
)

type LoginCredentials struct {
	Username string
	Password string
	GuardCode string
}

type SteamAuth struct {
	SessionId string
	SecureLogin string
	BrowserId string
}

func (sa *SteamAuth) Login(credentials *LoginCredentials) {
	
	loginInfo := fmt.Sprintf("username=%s&password=%s", url.QueryEscape(credentials.Username), url.QueryEscape(credentials.Password))
	
	req, _ := http.NewRequest("POST", "https://store.steampowered.com/login/DefaultAction", bytes.NewBuffer([]byte(loginInfo)))
	
	c := &http.Cookie{Name: "sessionid", Value:sa.SessionId}
	req.AddCookie(c)
	c = &http.Cookie{Name: "browserid", Value:sa.BrowserId}
	req.AddCookie(c)
	
	client := &http.Client{}
	
	resp, _ := client.Do(req)
	
	for _, cook := range resp.Cookies() {
		fmt.Printf("%s : %s", cook.Name, cook.Value)
	}
	
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	
	fmt.Println(string(body))
	
}

func (sa *SteamAuth) GetInitCookies() {
	
	loginUrl := "https://store.steampowered.com/login/"
	
	resp, err := http.Get(loginUrl)
	
	defer resp.Body.Close()
	
	if err != nil {
		fmt.Printf("Error HTTP : %v\n", err)
	}
	
	for _, c := range resp.Cookies() {
		fmt.Printf("%s: %s, %v\n", c.Name, c.Value, c.Expires)
		if c.Name == "browserid" {
			sa.BrowserId = c.Value
		} else if c.Name == "sessionid" {
			sa.SessionId = c.Value
		}
	}
}