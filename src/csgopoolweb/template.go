package csgopoolweb

import (
	"html/template"
	"io/ioutil"
	"fmt"
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
	
	mypoolLink := &Link{Caption: "My pool", Url: "/userpool/"}
	
	menu := fmt.Sprintf(`<ul><li>%s</li><li>2</li><li>3</li></ul>`, mypoolLink.GetHTML())
	
	return template.HTML(menu)
}