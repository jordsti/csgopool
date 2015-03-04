package csgopoolweb

import (
	"net/http"
	"fmt"
	"encoding/json"
	"csgodb"
)



func JSONHandler(w http.ResponseWriter, r *http.Request) {
	
	json_data := ""
	
	data := r.FormValue("data")
	
	db, _ := csgodb.Db.Open()
	
	if data == "ranking" {
		
		users := csgodb.GetUserPoints(db)
		
		b, _ := json.MarshalIndent(users, "", "	")
		json_data = string(b)
		
	} else if data == "matches" {
		matches := csgodb.GetAllMatches(db)
		b, _ := json.MarshalIndent(matches, "", "	")
		json_data = string(b)
	} else if data == "match" {
		m_id := ParseInt(r.FormValue("id"))
		match := csgodb.GetMatchById(db, m_id)
		match.FetchStats(db)
		b, _ := json.MarshalIndent(match, "", "	")
		json_data = string(b)
	} else {
		//default behaviour
		//send available options
	}
	
	db.Close()
	
	fmt.Printf("URL : %v\n", r.URL.Path[1:])
	
	fmt.Fprintf(w, json_data)
}