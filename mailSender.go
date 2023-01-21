package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type Mail struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func MailHTML(from string, to []string, password string, subject string, bodyMsg string) {

	t, _ := template.ParseFiles("templates/test.html")

	var tpl bytes.Buffer

	t.Execute(&tpl, struct {
		Title   string
		Message string
	}{
		Title:   "Title",
		Message: bodyMsg,
	})

	result := tpl.String()

	request := Mail{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    result,
	}

	msg := MessageBuilder(request)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func MessageBuilder(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg += fmt.Sprintf("From: %s\\n", mail.From)
	msg += fmt.Sprintf("To: %s\\n", mail.To)
	msg += fmt.Sprintf("Subject: %s\\n\\n", mail.Subject)
	msg += fmt.Sprintf("%s\\n", mail.Body)
	return msg
}

func MailWithAttachment() {
	//TODO
}
