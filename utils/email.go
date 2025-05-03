package utils

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmail отправляет код на почту
func SendEmail(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpUser == "" || smtpPass == "" {
		log.Println("ВНИМАНИЕ: SMTP не настроен. Используется эмуляция отправки email.")
		log.Printf("Эмуляция отправки email: To=%s, Subject=%s, Body=%s\n", to, subject, body)
		return nil // Возвращаем nil, чтобы продолжить процесс в тестовой среде
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		smtpHost,
		465,
		smtpUser,
		smtpPass,
	)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Ошибка отправки email: %v\n", err)
		return fmt.Errorf("ошибка отправки email: %w", err)
	}

	log.Printf("Email успешно отправлен: To=%s, Subject=%s\n", to, subject)
	return nil
}
