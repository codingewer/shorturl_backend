package auth

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"short-link/models"
)

func SendForgotPasswordEmail(email, token string) error {
	settings := models.Settings{}
	sitedata, _ := settings.FindBySiteName("short-url")
	from := sitedata.SmtpMail
	password := sitedata.SmtpPassword
	toList := []string{email}
	host := "smtp.gmail.com"
	port := "587"

	tmpl := template.Must(template.New("").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>{{ .Subject }}</title>
	</head>
	<body>
	<h1>{{ .Subject }}</h1>
	<p>{{ .Body }}</p>
	<a href="{{ .Link }}">{{ .Link }}</a>
	</body>
	</html>
	`))

	// HTML verisi
	data := struct {
		Subject string
		Body    string
		Link    string
	}{
		Subject: "Şifre Yenileme",
		Body:    "Link kısaltma uygulammıza tekrar giriş yapabilmek için 10 dk geçerli şifre yenileme linkiniz:",
		Link:    token,
	}

	// HTML'yi byte dizisine dönüştür
	var html bytes.Buffer
	if err := tmpl.Execute(&html, data); err != nil {
		panic(err)
	}

	// E-posta mesajı
	msg := []byte(
		"From: " + from + "\r\n" +
			"Subject: " + data.Subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=utf-8\r\n" +
			"\r\n" +
			html.String())
	//body := []byte(msg)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, toList, msg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}
