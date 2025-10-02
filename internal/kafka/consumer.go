package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"social-platform-kafka-worker/internal/model"
	"social-platform-kafka-worker/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader       *kafka.Reader
	emailService *service.EmailService
}

func NewConsumer(brokers, topic, groupID string, emailService *service.EmailService) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: strings.Split(brokers, ","),
			Topic:   topic,
			GroupID: groupID,
		}),
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
