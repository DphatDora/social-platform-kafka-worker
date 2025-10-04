package config

import (
	"log"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	once   sync.Once
	config Config
)

type Config struct {
	App      App
	Database Database
	Kafka    Kafka
	Email    Email
}

func LoadConfig() {
	once.Do(func() {
		// load .env
		_ = godotenv.Load()

		// load yaml config
		viper.SetConfigName("config") // config.yaml
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Config file error: %s", err)
		}

		// env override
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		bindEnvs()

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Config unmarshal error: %s", err)
		}
	})
}

func GetConfig() Config {
	LoadConfig()
	return config
}

func bindEnvs() {
	// Database
	_ = viper.BindEnv("database.url", "DB_URL")
	_ = viper.BindEnv("database.username", "DB_USER")
	_ = viper.BindEnv("database.password", "DB_PASSWORD")
	_ = viper.BindEnv("database.host", "DB_HOST")
	_ = viper.BindEnv("database.port", "DB_PORT")
	_ = viper.BindEnv("database.name", "DB_NAME")
	_ = viper.BindEnv("database.sslMode", "DB_SSLMODE")
	_ = viper.BindEnv("database.timeZone", "DB_TIMEZONE")

	// Kafka
	_ = viper.BindEnv("kafka.brokers", "KAFKA_HOST")
	_ = viper.BindEnv("kafka.topic", "KAFKA_TOPIC")
	_ = viper.BindEnv("kafka.groupId", "KAFKA_GROUP_ID")
	_ = viper.BindEnv("kafka.username", "KAFKA_USERNAME")
	_ = viper.BindEnv("kafka.password", "KAFKA_PASSWORD")
	_ = viper.BindEnv("kafka.securityProtocol", "KAFKA_SECURITY_PROTOCOL")
	_ = viper.BindEnv("kafka.saslMechanism", "KAFKA_SASL_MECHANISM")

	// Email
	_ = viper.BindEnv("email.user", "EMAIL_USER")
	_ = viper.BindEnv("email.password", "EMAIL_PASSWORD")
	_ = viper.BindEnv("email.smtpServer", "EMAIL_SMTP_SERVER")
	_ = viper.BindEnv("email.smtpPort", "EMAIL_SMTP_PORT")
}
