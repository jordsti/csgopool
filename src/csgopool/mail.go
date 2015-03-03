package csgopool

import (
	"net/smtp"
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
	return []byte(e.Message)
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
