package main

import (
	"context"
	"log"
	"time"

	"social-platform-kafka-worker/config"
	"social-platform-kafka-worker/internal/database"
	"social-platform-kafka-worker/internal/handler"
	"social-platform-kafka-worker/internal/kafka"
	"social-platform-kafka-worker/internal/repository"
	"social-platform-kafka-worker/internal/service"
)

const (
	DefaultPort = 8050
)

func main() {
	setUpInfrastructure()
	defer closeInfrastructure()
}

func setUpInfrastructure() {
	// Set time to UTC
	time.Local = time.UTC

	conf := config.GetConfig()
	log.Printf("[DEBUG] Config loaded: %+v", conf)

	// Init DB
	database.InitPostgresql(&conf)

	// Repo
	taskRepo := repository.NewTaskRepository(database.GetDB())

	// Service
	taskService := service.NewTaskService(taskRepo)
	emailService := service.NewEmailService(&conf)

	// Kafka
	producer := kafka.NewProducer(conf.Kafka)
	consumer := kafka.NewConsumer(conf.Kafka, emailService)

	// Handler
	taskHandler := handler.NewTaskHandler(taskService, producer)

	log.Printf("\n-----------------\n✅✅ Kafka Worker is running ✅✅\n-----------------")

	// Run Producer in background
	go func() {
		ctx := context.Background()
		for {
			taskHandler.ProcessDueTasks(ctx)
			time.Sleep(5 * time.Second)
		}
	}()

	// Run Consumer
	consumer.Start(context.Background())

	// port := conf.App.Port
	// if port == 0 {
	// 	port = DefaultPort
	// }
	// log.Printf("[✅] Kafka Worker running on port %d", port)
}

func closeInfrastructure() {
	if err := database.ClosePostgresql(); err != nil {
		log.Printf("[ERROR] Close postgresql fail: %s\n", err)
	}
}
