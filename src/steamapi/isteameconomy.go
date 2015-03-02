package steamapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"fmt"
)

const (
	USCurrency = 1
	EuroCurrency = 3
	
	CSGOAppId = 730
)

type ISteamEconomyResult struct {	
	Result map[string]*json.RawMessage `json:"result"`
	AssetInfos []AssetInfo
}

type MarketPrice struct {
	JSONLowestPrice string `json:"lowest_price"`
	JSONMedianPrice string `json:"median_price"`
	
	LowestPrice float32
	MedianPrice float32
	
	Volume string `json:"volume"`
}

type AssetInfo struct {
	AppId int
	IconUrl string `json:"icon_url"`
	IconUrlLarge string `json:"icon_url_large"`
	IconDragUrl string `json:"icon_drag_url"`
	Name string `json:"name"`
	MarketHashName string `json:"market_hash_name"`
	MarketName string `json:"market_name"`
	NameColor string `json:"name_color"`
	BackgroundColor string `json:"background_color"`
	Type string `json:"type"`
	Tradable string `json:"tradable"`
	Marketable string `json:"marketable"`
	Commodity string `json:"commodity"`
	FraudWarnings string `json:"fraudwarnings"`
	JSONDescriptions map[string]*json.RawMessage `json:"descriptions"`
	OwnerDescriptions string `json:"owner_descriptions"`
	JSONTags map[string]*json.RawMessage `json;"tags"`
	ClassId string `json:"classid"`
}

func (mp *MarketPrice) ParseJSON() {
		
	re := regexp.MustCompile(`([0-9]+\.[0-9]+)`)
	rs := re.FindStringSubmatch(mp.JSONLowestPrice)
	
	_price, _ := strconv.ParseFloat(rs[1], 32)
	
	mp.LowestPrice = float32(_price)
	
	rs = re.FindStringSubmatch(mp.JSONMedianPrice)
	
	_price, _ = strconv.ParseFloat(rs[1], 32)
	
	mp.MedianPrice = float32(_price)
}

func (a AssetInfo) GetPrice(currency int, country string) *MarketPrice {
	
	//building url
	url := fmt.Sprintf(`http://steamcommunity.com/market/priceoverview/?country=%s&currency=%d&appid=%d&market_hash_name=%s`,
		country, currency, a.AppId, url.QueryEscape(a.MarketHashName))
	
	resp, err := http.Get(url)
	
	if err != nil {
		fmt.Printf("HTTP Error : %v\n", err)
	}
	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		fmt.Printf("IO Error : %v\n", err)
	}

	price := &MarketPrice{}
	
	err = json.Unmarshal(body, price)
	
	if err != nil {
		fmt.Printf("Error [JSON]: %v\n", err)
	}
	
	price.ParseJSON()
	
	return price
}

func GetAssetClassInfo(appId int, classIds []string, key string) *ISteamEconomyResult {
	//building url
	url := fmt.Sprintf("https://api.steampowered.com/ISteamEconomy/GetAssetClassInfo/v1/?key=%s&appid=%d&class_count=%d",
			key, appId, len(classIds))
	
	//adding ids
	for i, id := range classIds {
		url += fmt.Sprintf("&classid%d=%s", i, id)
	}
	
	resp, err := http.Get(url)
	
	if err != nil {
		fmt.Printf("HTTP Error : %v\n", err)
	}
	
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		fmt.Printf("IO Error : %v\n", err)
	}
	
	results := &ISteamEconomyResult{}
	
	err = json.Unmarshal(body, results)
	
	if err != nil {
		fmt.Printf("Error [JSON]: %v\n", err)
	}
	
	for kn, dt := range results.Result {
		if kn != "success" {
			ai := AssetInfo{}
			ai.AppId = appId
			err = json.Unmarshal(*dt, &ai)
			
			if err != nil {
				fmt.Printf("Error [JSON]: %v\n", err)
			}
			
			results.AssetInfos = append(results.AssetInfos, ai)
		}
	}
	
	return results
}