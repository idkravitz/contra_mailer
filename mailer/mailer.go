package mailer

import "net/mail"

// Message is a fully described message to send
type Message struct {
	To   mail.Address
	From mail.Address

	Subject string
	Body    string
}

// Mailer provides a generic mail sending interface
type Mailer interface {
	Send(*Message) error
}
