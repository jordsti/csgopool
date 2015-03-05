package csgopoolweb

import (
	"net/http"
	"fmt"
	"encoding/json"
	"csgodb"
)

type JSONArgs struct {
	Name string
	Description string
}

type JSONFunc struct {
	Action string
	Description string
	Args []*JSONArgs
}

type JSONService struct {
	Actions []*JSONFunc
	Name string
	Version string
}

func GetFuncTable() *JSONService {
	service := &JSONService { Name:"CS:GO Pool", Version:"v0001"}
	
	funcs := []*JSONFunc{}
	
	f := &JSONFunc { Action: "ranking", Description: "Return the current User ranking with points" }
	funcs = append(funcs, f)
	
	f = &JSONFunc { Action: "matches", Description: "Return all matches played" }
	funcs = append(funcs, f)
	
	f = &JSONFunc { Action: "players", Description: "Return all players with their points" }
	funcs = append(funcs, f)
	
	f = &JSONFunc { Action: "match", Description: "Return match information and stats" }
	arg0 := &JSONArgs { Name:"id", Description:"Match Id" }
	f.Args = append(f.Args, arg0)
	funcs = append(funcs, f)
	
	service.Actions = funcs
	return service
	
}


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
	} else if data == "players" {
		players := csgodb.GetAllPlayersPoint(db)
		b, _ := json.MarshalIndent(players, "", "	")
		json_data = string(b)
	} else {
		funcs := GetFuncTable()
		b, _ := json.MarshalIndent(funcs, "", "	")
		json_data = string(b)
	}
	
	db.Close()
	
	fmt.Fprintf(w, json_data)
}