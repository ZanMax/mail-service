package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type Mail struct {
	From    string
	To      []string
	Subject string
	Body    string
}

func MailSender(from string, to []string, password string, subject string, bodyMsg string, attachments []string, templateFile string, params interface{}) {

	t, _ := template.ParseFiles(fmt.Sprintf("templates/%s", templateFile))

	var tpl bytes.Buffer

	errTemplate := t.Execute(&tpl, params)
	if errTemplate != nil {
		fmt.Println(errTemplate)
		return
	}

	request := Mail{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    tpl.String(),
	}

	var resultMessage []byte
	if len(attachments) > 0 {
		msg := MessageBuilderWithAttachment(request, attachments)
		resultMessage = []byte(msg)
	} else {
		msg := MessageBuilder(request)
		resultMessage = []byte(msg)
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, resultMessage)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func MessageBuilder(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.From)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func MessageBuilderWithAttachment(mail Mail, attachments []string) []byte {

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))

	boundary := "next-part"
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n",
		boundary))

	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	buf.WriteString(fmt.Sprintf("\r\n%s", mail.Body))
	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))

	for _, attachment := range attachments {
		buf.WriteString("Content-Transfer-Encoding: base64\r\n")
		buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", attachment))
		buf.WriteString("\r\n")

		data := readFile(attachment)

		b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(b, data)
		buf.Write(b)
		buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	}

	return buf.Bytes()
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
