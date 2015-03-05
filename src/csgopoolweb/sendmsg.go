package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgodb"
)

type WriteMessagePage struct {
	Page
	RecipientId string
	Recipient string
	Message string
}

func WriteMessageHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	if !session.IsLogged() {
		http.Redirect(w, r, "/", 302)
	}
	
	recipientId := ParseInt(r.FormValue("recipient_id"))
	action := r.FormValue("action")
	t, err := MakeTemplate("sendmsg.html")
	if err != nil {
		fmt.Println(err)
	}
	p := &WriteMessagePage{}
	
	m := GetMenu(session)
	db, _ := csgodb.Db.Open()
	recipient := csgodb.GetUserById(db, recipientId)
	
	if action == "send"  {
		p.Message = "Message sended !"
		
		title := r.FormValue("title")
		text := r.FormValue("text")
		
		title = template.HTMLEscapeString(title)
		text = template.HTMLEscapeString(text)
		
		csgodb.AddMessage(db, session.UserId, recipientId, title, text, csgodb.UnreadStatus)
		
	} else {
		if recipient.Id != 0 {
			p.Recipient = recipient.Name
			p.RecipientId = fmt.Sprintf("%d", recipient.Id)
		}
	}
	
	db.Close()

	p.Title = "CS:GO Pool - Inbox"
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	//p.LeftSide = template.HTML(curevent)
	p.GenerateRightSide(session)

	t.Execute(w, p)
	
}