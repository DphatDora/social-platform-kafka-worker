package config

type Kafka struct {
	Brokers          string
	Topic            string
	GroupID          string
	Username         string
	Password         string
	SecurityProtocol string
	SASLMechanism    string
}
