package csgopoolweb

import (
	"html/template"
	"net/http"
	"fmt"
	"csgodb"
)

type AccountFormPage struct {
	Page
	PasswordMinChar string
	UsernameMinChar string
	UsernameMaxChar string
}

func AccountFormHandler(w http.ResponseWriter, r *http.Request) {
	
	session := state.HandleSession(w, r)
	
	if session.IsLogged() {
		http.Redirect(w, r, "/", 301)
	}
	
	msgHtml := ""
	if session.IsFieldExists("message") {
		field := session.PopField("message")
		msgHtml = fmt.Sprintf(`<div>%s</div>`, field.Value)
	}
	
	t, err := MakeTemplate("accountform.html")
	if err != nil {
		state.Log.Error(fmt.Sprintf("%s", err))
	}
	
	constraints := csgodb.DefaultUserConstraints()

	p := &AccountFormPage{}
	p.PasswordMinChar = fmt.Sprintf("%d", constraints.PasswordMin)
	p.UsernameMinChar = fmt.Sprintf("%d", constraints.NameMin)
	p.UsernameMaxChar = fmt.Sprintf("%d", constraints.NameMax)
	p.Brand = "CS:GO Pool"
	p.Title = "CS:GO Pool - Create an account"
	p.Menu = template.HTML(GetMenu(session).GetHTML())
	p.Message = template.HTML(msgHtml)
	t.Execute(w, p)
}