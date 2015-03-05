package csgopoolweb

import (
	"fmt"
	"strings"
)

const(
	String = 1
	Int = 2
	Float = 3
)

type FormField interface {
	Name() string
	Value() string
	Type() int
}

type Field struct {
	name string
}

func (f *Field) SetName(name string) {
	f.name = name
}

func (f *Field) Name() string {
	return f.name
}

type StringField struct {
	Field
	value string
}

func (sf *StringField) SetValue(value string) {
	sf.value = value
}

func (sf *StringField) Value() string {
	return sf.value
}

func (sf *StringField) Type() int {
	return String
}

type IntField struct {
	Field
	value int
}

func (ifi *IntField) SetValue(value int) {
	ifi.value = value
}

func (ifi *IntField) Value() string {
	return fmt.Sprintf("%d", ifi.value)
}

func (ifi *IntField) Type() int {
	return Int
}

type FloatField struct {
	Field
	value float32
}

func (ff *FloatField) SetValue(value float32) {
	ff.value = value
}

func (ff *FloatField) Value() string {
	return fmt.Sprintf("%.2f", ff.value)
}

func (ff *FloatField) Type() int {
	return Float
}


func ReplaceFieldWithValue(content string, fields []FormField) string {
	tmp := content
	for _, ff := range fields {
		fmt.Printf("Field {%s, %s}\n", ff.Name(), ff.Value())
		tmp = strings.Replace(tmp, fmt.Sprintf("{{.%s}}", ff.Name()), ff.Value(), -1)
	}
	return tmp	
}