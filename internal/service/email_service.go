package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"social-platform-kafka-worker/config"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailService struct {
	SMTPHost string
	SMTPPort string
	User     string
	Password string
}

func NewEmailService(conf *config.Config) *EmailService {
	return &EmailService{
		SMTPHost: conf.Email.SMTPServer,
		SMTPPort: conf.Email.SMTPPort,
		User:     conf.Email.User,
		Password: conf.Email.Pass,
	}
}

func (s *EmailService) SendEmail(payload json.RawMessage) {
	var p EmailPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("❌ Invalid email payload: %v", err)
		return
	}

	auth := smtp.PlainAuth("", s.User, s.Password, s.SMTPHost)

	// log email content
	log.Printf("📧 Sending email to: %s, Subject: %s, Body: %s", p.To, p.Subject, p.Body)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Content-Transfer-Encoding: 7bit\r\n"+
		"\r\n%s",
		s.User, p.To, p.Subject, p.Body))

	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)
	err := smtp.SendMail(addr, auth, s.User, []string{p.To}, msg)
	if err != nil {
		log.Printf("❌ Failed to send email: %v", err)
	} else {
		log.Printf("✅ Email sent to %s", p.To)
	}
}
