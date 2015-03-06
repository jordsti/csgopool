package csgopoolweb

const (
	ErrorAlert = 1
	InfoAlert = 2
	WarningAlert = 3
)

type Alert struct {
	Title string
	Text string
	Type int
}

