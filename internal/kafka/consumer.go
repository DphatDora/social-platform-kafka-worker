package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"strings"

	"social-platform-kafka-worker/config"
	"social-platform-kafka-worker/internal/model"
	"social-platform-kafka-worker/internal/service"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Consumer struct {
	reader       *kafka.Reader
	emailService *service.EmailService
}

func NewConsumer(kafkaConfig config.Kafka, emailService *service.EmailService) *Consumer {
	readerConfig := kafka.ReaderConfig{
		Brokers: strings.Split(kafkaConfig.Brokers, ","),
		Topic:   kafkaConfig.Topic,
		GroupID: kafkaConfig.GroupID,
	}

	// Configure SASL/SSL if enabled
	if kafkaConfig.SecurityProtocol == "SASL_SSL" {
		mechanism := plain.Mechanism{
			Username: kafkaConfig.Username,
			Password: kafkaConfig.Password,
		}
		readerConfig.Dialer = &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		}
	}

	return &Consumer{
		reader:       kafka.NewReader(readerConfig),
		emailService: emailService,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Kafka read error: %v", err)
			continue
		}
		var task model.BotTask
		if err := json.Unmarshal(m.Value, &task); err != nil {
			log.Printf("❌ Unmarshal error: %v", err)
			continue
		}
		log.Printf("📩 Received task: %+v", task)

		switch task.Action {
		case "send_email":
			c.emailService.SendEmail(task.Payload)
		default:
			log.Printf("⚠️ Unknown action: %s", task.Action)
		}
	}
}
