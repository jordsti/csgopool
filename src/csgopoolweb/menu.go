package csgopoolweb



type MenuItem struct {
	MenuId int
	LinkName string
	Link string
	Active bool
}

type Menu struct {
	Items []MenuItem
}

func (m Menu) GetHTML() string {
	
	html := ""
	
	for _, i := range m.Items {
		
		html = html + "<li><a href=\""+i.Link+"\">"+i.LinkName+"</a></li>\n"
		
	}
	
	return html
}