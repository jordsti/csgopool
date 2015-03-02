package steamapi

import (
	"encoding/json"
	"github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/tradeoffer"
	"fmt"
	"io/ioutil"
)

type AccountCredentials struct {
	Username string
	Password string
	GuardCode string
	MachineHash []byte
}

type SteamState struct {
	ConfigPath string
	Connected bool
	Client *steam.Client
	Credentials *AccountCredentials
}

var Steam *SteamState

func NewClient(dataPath string) {
	Steam = &SteamState{}
	Steam.ConfigPath = dataPath
	Steam.Credentials = &AccountCredentials{}
	Steam.Credentials.Username = ""
	Steam.Credentials.Password = "not set"
	Steam.LoadConfig()
	
	if Steam.Credentials.Username == "" {
		fmt.Println("Invalid steam username")
	} else {
		Steam.Client = steam.NewClient()
	}
	
}

func (ss *SteamState) AcceptTrade(tradeId uint64) error {
	offerclient := tradeoffer.NewClient("", ss.Client.Web.SessionId, ss.Client.Web.SteamLogin, ss.Client.Web.SteamLoginSecure)
	return offerclient.Accept(tradeoffer.TradeOfferId(tradeId))
}

func (ss *SteamState) Connect() {
	
	if ss.Credentials.Username == "" {
		return
	}
	
	details := &steam.LogOnDetails{}
	details.Username = ss.Credentials.Username
	details.Password = ss.Credentials.Password
	
	if ss.Credentials.GuardCode != "" {
		details.AuthCode = ss.Credentials.GuardCode
	}
	

	details.SentryFileHash, _ = ioutil.ReadFile(ss.ConfigPath + "/steamhash")
	
	
	ss.Client.Connect()
	for event := range ss.Client.Events() {
		switch e := event.(type) {
		case *steam.ConnectedEvent:
			ss.Client.Auth.LogOn(details)
		case *steam.MachineAuthUpdateEvent:
			ioutil.WriteFile(ss.ConfigPath + "/steamhash" ,e.Hash, 0644)
			ss.SaveConfig()
		case *steam.LoggedOnEvent:
			ss.Connected = true
		case error:
			fmt.Println(e)
		}
		
		if ss.Connected && len(details.SentryFileHash) > 0 {
			break
		}
	}
	
}

func (ss *SteamState) SaveConfig() {
	filepath := ss.ConfigPath + "/steam.json"
	
	b, err := json.MarshalIndent(ss.Credentials, "", "	")
	
	if err != nil {
		fmt.Printf("JSON Error %v\n", err)
	}
	
	err = ioutil.WriteFile(filepath, b, 0644)
	
	if err != nil {
		fmt.Printf("IO Error %v\n", err)
	}
}

func (ss *SteamState) LoadConfig() {
	
	filepath := ss.ConfigPath + "/steam.json"
	
	b, err := ioutil.ReadFile(filepath)
	
	if err != nil {
		fmt.Printf("IO Error %v\n", err)
		ss.SaveConfig()
	}
	
	err = json.Unmarshal(b, ss.Credentials)
	
	if err != nil {
		fmt.Printf("JSON Error %v\n", err)
	}
}