package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/mail"
	"os"

	"github.com/kravitz/contra_lib/db"
	"github.com/kravitz/contra_lib/util"
	"github.com/kravitz/contra_mailer/mailer"
	"github.com/kravitz/contra_mailer/mailer/gmail"
	"github.com/streadway/amqp"
)

type mailerConfig struct {
	Homepage     string `json:"homepage"`
	SupportEmail string `json:"support_email"`
}

type greetData struct {
	Username     string
	Homepage     string
	SupportEmail string
}

type mailerApp struct {
	templatesPool map[string]*template.Template
	config        *mailerConfig
	mailSender    mailer.Mailer
	qCon          *amqp.Connection
}

func (app *mailerApp) renderTemplate(templateName string, data interface{}) (output string, err error) {
	tmpl, ok := app.templatesPool[templateName]
	if !ok {
		tmpl, err = template.ParseFiles("templates/" + templateName)
		if err != nil {
			return "", err
		}
		app.templatesPool[templateName] = tmpl
	}

	var htmlOut bytes.Buffer
	err = tmpl.Execute(&htmlOut, data)
	if err != nil {
		return "", err
	}
	bOutput := htmlOut.Bytes()

	return string(bOutput), err
}

func readConfig() (config *mailerConfig, err error) {
	inh, err := os.Open("./config.json")
	if err != nil {
		return nil, err
	}
	defer inh.Close()

	config = &mailerConfig{}
	configRaw, err := ioutil.ReadAll(inh) //inh.Read(configRaw)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configRaw, config)
	if err != nil {
		return nil, err
	}
	return config, err
}

func (app *mailerApp) sendGreetings(username string, email string) {
	body, err := app.renderTemplate("greet.tmpl", &greetData{
		Username:     username,
		Homepage:     app.config.Homepage,
		SupportEmail: app.config.SupportEmail,
	})
	if err != nil {
		log.Fatal(err)
	}
	msg := &mailer.Message{
		To:      mail.Address{Name: username, Address: email},
		From:    mail.Address{Name: "CONTRA service", Address: app.config.SupportEmail},
		Subject: "Welcome to CONTRA service, " + username,
		Body:    body,
	}
	err = app.mailSender.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rabbitUser := util.GetenvDefault("RABBIT_USER", "guest")
	rabbitPassword := util.GetenvDefault("RABBIT_PASSWORD", "guest")

	amqpSocket := fmt.Sprintf("amqp://%v:%v@tram-rabbit:5672", rabbitUser, rabbitPassword)
	amqpCon, err := db.RabbitInitConnect(amqpSocket)

	if err != nil {
		log.Fatal(err)
	}

	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	m, err := gmail.CreateMailer("client_secret.json", "credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	app := &mailerApp{
		templatesPool: map[string]*template.Template{},
		config:        config,
		mailSender:    m,
		qCon:          amqpCon,
	}

	app.sendGreetings("idkravitz", "idkravitz@gmail.com")
	fmt.Printf("Sent grettings to idkravitz!")
}
