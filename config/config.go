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
		viper.SetEnvPrefix("SOCIALPLATFORM")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("Config unmarshal error: %s", err)
		}
	})
}

func GetConfig() Config {
	LoadConfig()
	return config
}
