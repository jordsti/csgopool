package csgopoolweb

import (
	//"html/template"
	"net/http"
	"fmt"
	"csgodb"
)


func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	
	username := r.FormValue("username")
	pwd := r.FormValue("password")
	pwd2 := r.FormValue("password2")
	email := r.FormValue("email")
	
	session := state.HandleSession(w, r)
	
	if session.IsLogged() {
		http.Redirect(w, r, "/", 301)
		return
	}
	
	db, _ := csgodb.Db.Open()
	
	if pwd == pwd2 {
		err := csgodb.CreateUser(db, username, pwd, email, csgodb.UserRank)
		
		if err != nil {
			session.SetField("message", fmt.Sprintf("%s", err))
		} else {
			state.Log.Info(fmt.Sprintf("User %s created", username))
			session.SetField("message", "Account created with success")
		}
	} else {
		session.SetField("message", "Password mismatch")
	}
	
	http.Redirect(w, r, "/accountform/", 301)
}