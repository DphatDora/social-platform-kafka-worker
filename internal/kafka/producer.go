package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
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
