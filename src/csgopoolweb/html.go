package csgopoolweb

import (
	"fmt"
)

type HtmlAttr struct {
	Name string
	Value string
}

type HtmlElement struct {
	Tag string
	Attributes []*HtmlAttr
	Childs []*HtmlElement
	InnerText string
}

func (he *HtmlElement) SetAttribute(name string, value string) {
	attr := &HtmlAttr{Name: name, Value: value}
	he.Attributes = append(he.Attributes, attr)
}

func (he *HtmlElement) AddChild(child *HtmlElement) {
	he.Childs = append(he.Childs, child)
}

func (he *HtmlElement) GetHTML() string {
	
	attrs := ""
	
	for _, attr := range he.Attributes {
		if len(attr.Value) > 0 {
			attrs += fmt.Sprintf(`%s="%s" `, attr.Name, attr.Value)
		} else {
			attrs += fmt.Sprintf("%s ", attr.Name)
		}
	}
	
	childs := ""
	
	for _, child := range he.Childs {
		childs += child.GetHTML()
	}
	
	html := fmt.Sprintf(`<%s %s>
	%s
	%s
	 </%s>`, he.Tag, attrs, childs, he.InnerText, he.Tag)
	
	return html
}

func CreateDiv() *HtmlElement {
	element := &HtmlElement{Tag: "div"}
	return element
}