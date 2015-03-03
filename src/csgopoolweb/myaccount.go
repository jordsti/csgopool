package csgopoolweb

import (
	"html/template"
	"net/http"
	"csgodb"
	"csgopool"
	"strconv"
	"fmt"
)

type MyAccountPage struct {
	Page
	SteamID string
	Credit string
	Email string
	Transactions template.HTML
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
	transactions := csgodb.GetTransactionsByUser(db, session.UserId)
	
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
	
	transactions_html := ""
	
	for _, ts := range transactions {
		
		str_amount := ""
		if ts.Amount > 0 {
			str_amount = fmt.Sprintf("+%.2f", ts.Amount)
		} else {
			str_amount = fmt.Sprintf("%.2f", ts.Amount)
		}
		
		tsLink := &Link{Caption: "View", Url:"/viewtransaction/"}
		tsLink.AddInt("id", ts.TransactionId)
		
		date_str := fmt.Sprintf("%d-%02d-%02d %02d:%02d", 
			ts.Timestamp.Year(), ts.Timestamp.Month(), ts.Timestamp.Day(), ts.Timestamp.Hour(), ts.Timestamp.Minute())
		
		transactions_html += fmt.Sprintf(`
		<tr>
			<td>%d</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>
		`, ts.TransactionId, ts.Description, str_amount, date_str, tsLink.GetHTML())
		
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
	p.Version = csgopool.CurrentVersion.String()
	p.Transactions = template.HTML(transactions_html)
	p.Email = session.User.Email
	t.Execute(w, p)
}