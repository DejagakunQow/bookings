package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"bookings/cmd/web/internal/models"

	mail "github.com/xhit/go-simple-mail"
)

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

	// CORRECT TEMPLATE LOGIC
	if m.Template == "" {
		// No template → send content directly
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		// Load template file
		data, err := ioutil.ReadFile(fmt.Sprintf("../email-templates/%s", m.Template))
		if err != nil {
			errorLog.Println("Template read error:", err)
			return
		}

		templateStr := string(data)

		// Replace placeholder with your content
		msgToSend := strings.Replace(templateStr, "[%body%]", m.Content, 1)

		email.SetBody(mail.TextHTML, msgToSend)
	}

	// Send email
	err = email.Send(client)
	if err != nil {
		errorLog.Println("Email send error:", err)
	} else {
		log.Println("Email sent successfully.")
	}
}
