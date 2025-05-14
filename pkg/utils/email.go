package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends an email to the specified address
func SendEmail(to, subject, body string) error {
	// Get email configuration from environment
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	if smtpUsername == "" {
		log.Println("WARNING: SMTP_USERNAME not configured")
		// In development mode, log the email instead of sending it
		log.Printf("Would send email to %s\nSubject: %s\nBody: %s", to, subject, body)
		return nil
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		log.Println("WARNING: SMTP_PASSWORD not configured")
		// In development mode, log the email instead of sending it
		log.Printf("Would send email to %s\nSubject: %s\nBody: %s", to, subject, body)
		return nil
	}

	// Set up authentication for SMTP server
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Format the email
	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", smtpUsername, to, subject, body))

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUsername, []string{to}, message)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return fmt.Errorf("error sending email: %w", err)
	}

	log.Printf("Email sent to %s", to)
	return nil
}
