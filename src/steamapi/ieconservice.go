package steamapi

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	//account type
	IndividualType = 1
	
	//trade state
	InvalidTradeState = 1
	ActiveTradeState = 2
	AcceptedTradeState = 3
	CounteredTradeState = 4
	ExpiredTradeState = 5
	CanceledTradeState = 6
	DeclinedTradeState = 7
	InvalidItemsState = 8
	EmailCanceledState = 10
)

type CEcon_Asset struct {
	AppId string `json:"appid"`
	ContextId string `json:"contextid"`
	AssetId string `json:"assetid"`
	ClassId string `json:"classid"`
	InstanceId string `json:"instanceid"`
	Amount string `json:"amount"`
	Missing bool `json:"missing"`
}

type CEcon_TradeOffer struct {
	
	TradeOfferId string `json:"tradeofferid"`
	AccountIdOther int64 `json:"accountid_other"`
	Message string `json:"message"`
	ExpirationTime int64 `json:"expiration_time"`
	TradeOfferState int `json:"trade_offer_state"`
	ItemsToReceive []*CEcon_Asset `json:"items_to_receive"`
	ItemsToGive []*CEcon_Asset `json:"items_to_give"`
	IsOurOffer bool `json:"is_our_offer"`
	TimeCreated int64 `json:"time_created"`
	TimeUpdated int64 `json:"time_updated"`
	FromRealTimeTrade bool `json:"from_real_time_trade"`
}

type IConServiceResponse struct {
	TradeOffersReceived []*CEcon_TradeOffer `json:"trade_offers_received"`
}


type JSONResponse struct {
	Response IConServiceResponse `json:"response"`
}

func (to *CEcon_TradeOffer) SteamID() *SteamID {

	id := &SteamID{AccountId:to.AccountIdOther, Type: IndividualType}
	return id
}

func (r *JSONResponse) Parse(data []byte) {
	err := json.Unmarshal(data, r)

	if err != nil {
		fmt.Printf("Error [3]: %v\n", err)
	}
}

func GetTradeOffers(key string) IConServiceResponse {
	
	url := fmt.Sprintf(`https://api.steampowered.com/IEconService/GetTradeOffers/v1/?key=%s&format=json&input_json={"get_received_offers":true,"active_only":true}`, key)
	
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error [1] : %v\n", err)
	}
	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		fmt.Printf("Error [2]: %v\n", err)
	}
	
	response := &JSONResponse{}
	response.Parse(body)
	return response.Response
}

