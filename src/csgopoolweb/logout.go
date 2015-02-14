package csgopoolweb

import (
	"net/http"
)


func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	

	session := state.HandleSession(w, r)
	session.UserId = 0
	session.ClearFields()
	
	
	http.Redirect(w, r, "/", 301)
}