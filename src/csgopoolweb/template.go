package csgopoolweb

import (
	"html/template"
	"io/ioutil"
)

func MakeTemplate(contentPath string) (*template.Template, error) {
	
	html := ""
	
	headerFile := state.RootPath + "header.html"
	footerFile := state.RootPath + "footer.html"
	contentPath = state.RootPath + contentPath
	
	header, _ := ioutil.ReadFile(headerFile)
	footer, _ := ioutil.ReadFile(footerFile)
	content, _ := ioutil.ReadFile(contentPath)
	
	html = string(header) + string(content) + string(footer)
	
	t, err := template.New(contentPath).Parse(html)
	
	return t, err
	
}

func GetLoginForm() template.HTML {
	
	loginForm, _ := ioutil.ReadFile(state.RootPath + "login.html")
	loginHtml := string(loginForm)
	
	return template.HTML(loginHtml)
}

func GetUserMenu() template.HTML {
	
	menu := `<a href="/userpool/" class="btn btn-default">My Pool</a><br />`
	menu += `<a href="/myaccount/" class="btn btn-default">My Account</a><br />`
	menu += `<a href="/logout/" class="btn btn-default">Log Out</a><br />`
	return template.HTML(menu)
}