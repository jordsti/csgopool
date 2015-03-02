package main

import (
	"steamapi"
	"fmt"
)

func main() {
	
	
	response := steamapi.GetTradeOffers("")
	fmt.Printf("%d\n", len(response.TradeOffers()))
}

