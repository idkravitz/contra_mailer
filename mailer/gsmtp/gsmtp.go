package gsmtp

import (
	"encoding/base64"
	"fmt"
	"net/smtp"

	"github.com/kravitz/contra_mailer/mailer"
)

// Mailer provides generic SMTP mailing
type Mailer struct {
	auth    smtp.Auth
	address string
}

// CreateMailer is a constructor for gsmtp.Mailer
func CreateMailer(username, password, host, port string) (m *Mailer) {
	m = &Mailer{
		auth:    smtp.PlainAuth("", username, password, host),
		address: fmt.Sprintf("%v:%v", host, port),
	}
	return m
}

func encodeWeb64String(b []byte) string {
	s := base64.URLEncoding.EncodeToString(b)

	var i = len(s) - 1
	for s[i] == '=' {
		i--
	}

	return s[0 : i+1]
}

// Send mail
func (m *Mailer) Send(msg *mailer.Message) error {
	header := map[string]string{}
	header["From"] = msg.From.String()
	header["To"] = msg.To.String()
	header["Subject"] = msg.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	//header["Content-Transfer-Encoding"] = "base64"

	var rawMsg string
	for k, v := range header {
		rawMsg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	rawMsg += "\r\n" + msg.Body

	err := smtp.SendMail(m.address, m.auth, msg.From.Name, []string{msg.To.Address}, []byte(rawMsg))
	return err
}
