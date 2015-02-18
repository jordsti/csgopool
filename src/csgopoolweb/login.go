package csgopoolweb

import (
	//"html/template"
	"net/http"
	"csgodb"
)


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	
	username := r.FormValue("username")
	pwd := r.FormValue("password")
	
	session := state.HandleSession(w, r)
	
	if session.IsLogged() {
		http.Redirect(w, r, "/", 301)
		return
	}
	
	db, _ := csgodb.Db.Open()
	
	user, err := csgodb.Login(db, username, pwd)
	
	if err != nil {
		//fmt.Println("Login error")
		session.SetField("message", "Bad username/password combination")
	} else {
		//fmt.Printf("Login success [%s]\n", user.Name)
		session.UserId = user.Id
	}
	
	db.Close()
	
	http.Redirect(w, r, "/", 301)
}