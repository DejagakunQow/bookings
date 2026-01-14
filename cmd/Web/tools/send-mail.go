//go:build tools

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
	mail "github.com/xhit/go-simple-mail"
)

// dummy loggers for tool build
var errorLog = log.New(log.Writer(), "ERROR\t", log.LstdFlags)

// dummy app struct for tool build
type application struct {
	MailChan chan models.MailData
}

var app application

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println("SMTP connection error:", err)
		return
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).
		AddTo(m.To).
		SetSubject(m.Subject)

	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("../email-templates/%s", m.Template))
		if err != nil {
			errorLog.Println("Template read error:", err)
			return
		}

		templateStr := string(data)
		msgToSend := strings.Replace(templateStr, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	if err := email.Send(client); err != nil {
		errorLog.Println("Email send error:", err)
	} else {
		log.Println("Email sent successfully.")
	}
}

func main() {}
