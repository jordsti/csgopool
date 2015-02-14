package csgopoolweb

import (
	//"html/template"
	"net/http"
	"fmt"
	"csgopool"
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
	
	if pwd == pwd2 {
		u, err := state.Users.CreateUser(username, pwd, email, csgopool.UserRank)
		
		if err != nil {
			session.SetField("message", fmt.Sprintf("%s", err))
		} else {
			state.Log.Info(fmt.Sprintf("User %s created", u.Name))
			session.SetField("message", "Account created with success")
		}
	} else {
		session.SetField("message", "Password mismatch")
	}
	
	http.Redirect(w, r, "/accountform/", 301)
}