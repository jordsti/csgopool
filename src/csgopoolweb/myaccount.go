package csgopoolweb

import (
	"html/template"
	"net/http"
	"csgodb"
	"fmt"
)

type MyAccountPage struct {
	Page
	Email string
}

func MyAccountHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	if !session.IsLogged() {
		http.Redirect(w, r, "/", 302)
	}
	
	
	action := r.FormValue("action")
	if action == "chpwd" {
		//changing password
		curpass := r.FormValue("curpassword")
		password := r.FormValue("password")
		password2 := r.FormValue("password")
		

		if password != password2 {
			session.SetField("message", "Password mismatch")
		} else {
			
			db, _ := csgodb.Db.Open()
			
			_, err := csgodb.Login(db, session.User.Name, curpass)
			
			if err == nil {
				
				err = csgodb.UpdatePassword(db, session.UserId, password)				
				
				if err != nil {
					session.SetField("message", fmt.Sprintf("%s", err))
				} else {
					session.SetField("message", "Password changed with success!")
				}

			} else {
				session.SetField("message", "This is not your current password")
			}
			
			db.Close()
			
		}
		
	} 
	
	msgHtml := ""
	if session.IsFieldExists("message") {
		field := session.PopField("message")
		msgHtml = fmt.Sprintf(`<div>%s</div>`, field.Value)
	}
	
	t, err := MakeTemplate("myaccount.html")
	if err != nil {
		state.Log.Error(fmt.Sprintf("%s", err))
	}
	
	p := &MyAccountPage{}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - My Account"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Message = template.HTML(msgHtml)
	p.Email = session.User.Email
	t.Execute(w, p)
}