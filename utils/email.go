package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

// SendEmail отправляет код на почту
func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		465,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("ошибка отправки email: %w", err)
	}
	return nil
}
