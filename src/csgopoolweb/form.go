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

func (f Field) Name() string {
	return f.name
}

type StringField struct {
	Field
	value string
}

func (sf StringField) Value() string {
	return sf.value
}

func (sf StringField) Type() int {
	return String
}

type IntField struct {
	Field
	value int
}

func (ifi IntField) Value() string {
	return fmt.Sprintf("%d", ifi.value)
}

func (ifi IntField) Type() int {
	return Int
}

type FloatField struct {
	Field
	value float32
}

func (ff FloatField) Value() string {
	return fmt.Sprintf("%.2f", ff.value)
}

func (ff FloatField) Type() int {
	return Float
}


func ReplaceFieldWithValue(content string, fields []FormField) string {
	for _, ff := range fields {
		content = strings.Replace(content, fmt.Sprintf("{{.%s}}", ff.Name()), ff.Value(), -1)
	}
	return content	
}