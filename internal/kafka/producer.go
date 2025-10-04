package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"

	"social-platform-kafka-worker/config"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(kafkaConfig config.Kafka) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaConfig.Brokers),
		Topic:    kafkaConfig.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	// Configure SASL/SSL if enabled
	if kafkaConfig.SecurityProtocol == "SASL_SSL" {
		mechanism := plain.Mechanism{
			Username: kafkaConfig.Username,
			Password: kafkaConfig.Password,
		}
		writer.Transport = &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		}
	}

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) SendMessage(ctx context.Context, v interface{}) error {
	value, _ := json.Marshal(v)
	err := p.writer.WriteMessages(ctx, kafka.Message{Value: value})
	if err != nil {
		log.Printf("❌ Kafka send error: %v", err)
		return err
	}
	log.Println("✅ Sent message to Kafka")
	return nil
}
