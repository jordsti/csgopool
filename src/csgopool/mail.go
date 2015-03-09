package csgopool

import (
	"net/smtp"
	"fmt"
)

const (
	EmailPending = 0
	EmailSended = 1
	EmailError = 2
)

type Email struct {
	Destination []string
	From string
	Subject string
	Message string
	Status int
}

func (e *Email) Body() []byte {
	//todo
	//construct mail here
	email := ""
	email += fmt.Sprintf("From: %s\r\n", e.From)
	
	to_str := ""
	
	for _, d := range e.Destination {
		to_str += d + ", "
	}
	
	email += fmt.Sprintf("To: %s\r\n", to_str)
	email += fmt.Sprintf("Subject: %s\r\n", e.Subject)
	email += `Content-type: text/plain; charset="utf-8"\r\n`
	email += `MIME-Version: 1.0\r\n`
	//email += `Content-Transfer-Encoding: base64\r\n`
	email = fmt.Sprintf("%s\r\n%s", email, e.Message)

	return []byte(email)
}


func (e *Email) Send() {
	
	if Pool.Settings.Mail.Address != "" {
		auth := smtp.PlainAuth("", Pool.Settings.Mail.Username, Pool.Settings.Mail.Password, Pool.Settings.Mail.Host())
		err := smtp.SendMail(Pool.Settings.Mail.Host(), auth, e.From, e.Destination, e.Body())
		
		if err != nil {
			e.Status = EmailSended
		} else {
			e.Status = EmailError
		}
	}
}
