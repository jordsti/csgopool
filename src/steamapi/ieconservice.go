package steamapi

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type CEcon_Asset struct {
	appid int
	contextid int
	assetid int64
	classid int64
	instanceid int
	amount int
	missing bool
}

type CEcon_TradeOffer struct {
	
	tradeofferid int64
	accountid_other int64
	message string
	expiration_time int64
	trade_offer_state int
	items_to_receive []*CEcon_Asset
	items_to_give []*CEcon_Asset
	is_our_offer bool
	time_created int64
	time_updated int64
	from_real_time_trade bool
}

type IConServiceResponse struct {
	trade_offers_received []*CEcon_TradeOffer
}

func(r *IConServiceResponse) TradeOffers() []*CEcon_TradeOffer {
	return r.trade_offers_received
}

type JSONResponse struct {
	response IConServiceResponse
}


func (r *JSONResponse) Parse(data []byte) {
	err := json.Unmarshal(data, &r.response)

	if err != nil {
		fmt.Printf("Error [3]: %v\n", err)
	}
}

func GetTradeOffers(key string) IConServiceResponse {
	
	url := fmt.Sprintf(`https://api.steampowered.com/IEconService/GetTradeOffers/v1/?key=%s&format=json&input_json={"get_received_offers":true}`, key)
	
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
	
	fmt.Printf("%d", len(response.response.TradeOffers()))
	return response.response
}

