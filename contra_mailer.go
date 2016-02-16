package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

type greetData struct {
	Username     string
	Homepage     string
	SupportEmail string
}

type mailerApp struct {
	templatesPool map[string]*template.Template
}

func (app *mailerApp) renderTemplate(templateName string, data interface{}) (output []byte, err error) {
	tmpl, ok := app.templatesPool[templateName]
	if !ok {
		tmpl, err = template.ParseFiles("templates/" + templateName)
		if err != nil {
			return nil, err
		}
		app.templatesPool[templateName] = tmpl
	}

	var htmlOut bytes.Buffer
	err = tmpl.Execute(&htmlOut, data)
	if err != nil {
		return nil, err
	}
	output = htmlOut.Bytes()

	return output, err
}

func main() {
	app := &mailerApp{templatesPool: map[string]*template.Template{}}
	output, err := app.renderTemplate("greet.tmpl", &greetData{
		Username:     "idkravitz",
		Homepage:     "http://cuba.dvfu.ru/contra",
		SupportEmail: "cubadvfu@gmail.com",
	})
	if err != nil {
		log.Fatalf("Template parse failed: %v", err)
	}
	fmt.Print(string(output))
}
