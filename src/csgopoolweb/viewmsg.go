package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgodb"
)

type ViewMessagePage struct {
	Page
	Sender template.HTML
	SentOn string
	MessageTitle string
	Text string
	Links template.HTML
}

func ViewMessageHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	if !session.IsLogged() {
		http.Redirect(w, r, "/", 302)
	}

	t, err := MakeTemplate("viewmsg.html")
	if err != nil {
		fmt.Println(err)
	}
	
	
	m := GetMenu(session)

	db, _ := csgodb.Db.Open()
	
	m_id := ParseInt(r.FormValue("id"))
	
	message := csgodb.GetMessageById(db, m_id)
	p := &ViewMessagePage{}
	if message.SenderId == session.UserId || message.RecipientId == session.UserId {
		p.MessageTitle = message.Title
		p.Text = message.Text
		
		senderLink := &Link{Caption: message.SenderName, Url:"/viewuser/"}
		senderLink.AddInt("id", message.SenderId)
		
		p.Sender = template.HTML(senderLink.GetHTML())
		
		link := &Link{Caption: "Reply", Url:"/sendmsg/"}
		link.AddInt("recipient_id", message.SenderId)
		
		deleteLink := &Link{Caption: "Delete", Url:"/inbox/"}
		deleteLink.AddParameter("action", "delete")
		deleteLink.AddInt("id", message.MessageId)
		
		p.Links = template.HTML(fmt.Sprintf("%s | %s", link.GetHTML(), deleteLink.GetHTML()))
		
		//update message status
		
		message.UpdateStatus(db, csgodb.ReadedStatus)
		
	} else {
		http.Redirect(w, r, "/", 302)
	}
	
	
	db.Close()
	
	p.Title = fmt.Sprintf("CS:GO Pool - View Message : %s", message.Title)
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	//p.LeftSide = template.HTML(curevent)
	p.GenerateRightSide(session)
	p.SentOn = fmt.Sprintf("%d-%02d-%02d %02d:%02d", 
		message.SentOn.Year(), 
		message.SentOn.Month(), 
		message.SentOn.Day(), 
		message.SentOn.Hour(), 
		message.SentOn.Minute())
	
	t.Execute(w, p)
	
}