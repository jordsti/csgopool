package csgopoolweb

import (
	"net/http"
	"html/template"
	"fmt"
	"csgodb"
)

type InboxPage struct {
	Page
	Messages template.HTML
}

func InboxHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)

	if !session.IsLogged() {
		http.Redirect(w, r, "/", 302)
	}

	t, err := MakeTemplate("inbox.html")
	if err != nil {
		fmt.Println(err)
	}
	
	start := 0
	count := 50
	
	m := GetMenu(session)
	messages_html := ""
	db, _ := csgodb.Db.Open()
	
	messages := csgodb.GetUserMessages(db, session.UserId, start, start + count)
	
	for _, msg := range messages {
		msgLink := &Link{Caption: fmt.Sprintf("%d", msg.MessageId), Url:"/viewmsg/"}
		msgLink.AddInt("id", msg.MessageId)
		
		messages_html += fmt.Sprintf(`
		<tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>
		`, msgLink.GetHTML(), msg.SenderName, msg.Title, "")
	}
	
	db.Close()
	
	p := &InboxPage{}
	p.Title = "CS:GO Pool - Inbox"
	p.Brand = "CS:GO Pool"
	p.Menu = template.HTML(m.GetHTML())
	//p.LeftSide = template.HTML(curevent)
	p.GenerateRightSide(session)
	p.Messages = template.HTML(messages_html)
	t.Execute(w, p)
	
}