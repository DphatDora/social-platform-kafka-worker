package service

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"social-platform-kafka-worker/config"
	"social-platform-kafka-worker/package/logger"
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
		Password: conf.Email.Password,
	}
}

func (s *EmailService) SendEmail(payload json.RawMessage) {
	var p EmailPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		logger.Errorf("[Error] Invalid email payload: %v", err)
		return
	}

	auth := smtp.PlainAuth("", s.User, s.Password, s.SMTPHost)

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
		logger.Errorf("[Error] Failed to send email: %v", err)
	} else {
		logger.Infof("[Email] Email sent to %s", p.To)
	}
}
