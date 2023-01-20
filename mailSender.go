package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

func Mail(from string, to []string, password string, subject string, bodyMsg string) {

	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, _ := template.ParseFiles("templates/test.html")

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("To: %s\\n%s\\n\\n", subject, mimeHeaders)))

	t.Execute(&body, struct {
		Title   string
		Message string
	}{
		Title:   "Title",
		Message: bodyMsg,
	})
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
}
