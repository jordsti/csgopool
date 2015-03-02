package csgopoolweb

import (
	"html/template"
	"net/http"
	"csgodb"
	"strconv"
	"fmt"
)

type MyAccountPage struct {
	Page
	SteamID string
	Credit string
	Email string
}

func MyAccountHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	if !session.IsLogged() {
		http.Redirect(w, r, "/", 302)
	}
	
	db, _ := csgodb.Db.Open()
	action := r.FormValue("action")
	if action == "chpwd" {
		//changing password
		curpass := r.FormValue("curpassword")
		password := r.FormValue("password")
		password2 := r.FormValue("password")
		

		if password != password2 {
			session.SetField("message", "Password mismatch")
		} else {
			
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
		}
		
	} else if action == "steamid" {
		sid := r.FormValue("steamid")
		steamid := csgodb.GetSteamIDByUser(db, session.UserId)
		if len(sid) > 0 {
			_sid, _ := strconv.ParseInt(sid, 10, 64)
			
			if steamid.LinkId == 0 {
				csgodb.SaveSteamID(db, session.UserId, int64(_sid))
			} else {
				steamid.SteamId = int64(_sid)
				steamid.UpdateSteamID(db)
			}
		} 
	}
	
	steamid := csgodb.GetSteamIDByUser(db, session.UserId)
	credit := csgodb.GetCreditByUser(db, session.UserId)
	
	db.Close()
	
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
	
	if steamid.LinkId == 0 {
		p.SteamID = "NOT SET"
	} else {
		p.SteamID = fmt.Sprintf("%d", steamid.SteamId)
	}
	
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - My Account"
	p.Credit = fmt.Sprintf(`%.2f`, credit.Amount)
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Message = template.HTML(msgHtml)
	p.Email = session.User.Email
	t.Execute(w, p)
}