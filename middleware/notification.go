package middleware

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

func SendMail(toMail, subject, body string) error {
	// create new message
	message := gomail.NewMessage()

	// set email headers
	message.SetHeader("From", "nadyafa795@gmail.com")
	message.SetHeader("To", toMail)
	message.SetHeader("Subject", subject)

	// generate email body
	message.SetBody("text/plain", body)

	// get smtp setup
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// setup smtp dialer
	dialer := gomail.NewDialer(host, int(port), username, password)

	// tls config
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// send email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println(port)
		return fmt.Errorf("failed to send email. host: %s, port: %d, username: %s, password: %s, error: %v", host, port, username, password, err)
		// panic(err)
	}

	// fmt.Println("Email sent successfully")
	return nil
}
