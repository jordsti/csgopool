package csgopool

import (
	"steamapi"
	"csgodb"
	"fmt"
	"strconv"
)

func (w *WatcherState) WatchIncomingTrades() {
	
	key := Pool.Settings.SteamKey
	
	db, _ := csgodb.Db.Open()
	
	if len(key) == 0 {
		w.Log.Info("Invalid Steam Key !")
	} else {
		
		resp := steamapi.GetTradeOffers(key)
		
		for _, trade := range resp.TradeOffersReceived {
			sId := trade.SteamID()

			if len(trade.ItemsToGive) > 0 {
				w.Log.Info(fmt.Sprintf("Trade [%s] Invalid need to give items", trade.TradeOfferId))
				continue
			}
			
			if trade.TradeOfferState == steamapi.ActiveTradeState {
				user := csgodb.GetUserBySteamID(db, sId.AccountId)
				if user != nil {
					//accept this trade and add lowest value to credit
					tradeId, _ := strconv.ParseUint(trade.TradeOfferId, 10, 64)
					total := float32(0.00)
					//computing trade sums
					for _, i := range trade.ItemsToReceive {
						classid := []string{i.ClassId}
						rs := steamapi.GetAssetClassInfo(steamapi.CSGOAppId, classid, Pool.Settings.SteamKey)
						
						for _, ai := range rs.AssetInfos {
							p := ai.GetPrice(steamapi.USCurrency, "US")
							total += p.LowestPrice
						}
					}
					
					err := steamapi.Steam.AcceptTrade(tradeId)
					
					if err == nil {
						//trade completed
						//need to add credit to that account
						w.Log.Info(fmt.Sprintf("Adding %.2f credit to user [%d]", total, user.Id))
						
						credit := csgodb.GetCreditByUser(db, user.Id)
						
						if credit.CreditId == 0 {
							//new credit
							csgodb.InsertCredit(db, user.Id, total)
						} else {
							credit.Add(total)
							credit.UpdateCredit(db)
						}
					}
					
				} else {
					w.Log.Error(fmt.Sprintf("User not found for steam id : %d, skipping this trade", sId.AccountId))
				}	
			} else {
				w.Log.Info(fmt.Sprintf("Trade [%d] Invalid State", trade.TradeOfferId))
			}
			
			
		}
		
	}
	
	db.Close()
	
}

