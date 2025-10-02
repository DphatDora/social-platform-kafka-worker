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
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s",
		p.To, p.Subject, p.Body))

	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)
	err := smtp.SendMail(addr, auth, s.User, []string{p.To}, msg)
	if err != nil {
		log.Printf("❌ Failed to send email: %v", err)
	} else {
		log.Printf("✅ Email sent to %s", p.To)
	}
}
