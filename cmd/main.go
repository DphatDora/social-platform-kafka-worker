package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"social-platform-kafka-worker/config"
	"social-platform-kafka-worker/internal/database"
	"social-platform-kafka-worker/internal/handler"
	"social-platform-kafka-worker/internal/kafka"
	"social-platform-kafka-worker/internal/repository"
	"social-platform-kafka-worker/internal/service"
	"social-platform-kafka-worker/package/logger"
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
	if err := logger.Init(&conf.Log); err != nil {
		fmt.Fprintf(os.Stderr, "logger init: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = logger.Sync() }()

	logger.Debugf("[DEBUG] Config loaded: %+v", conf)

	// Init DB
	database.InitPostgresql(&conf)

	// Repo
	taskRepo := repository.NewTaskRepository(database.GetDB())
	userBadgeRepo := repository.NewUserBadgeRepository(database.GetDB())
	userRepo := repository.NewUserRepository(database.GetDB())
	interestScoreRepo := repository.NewInterestScoreRepository(database.GetDB())
	tagPreferenceRepo := repository.NewTagPreferenceRepository(database.GetDB())

	// Service
	taskService := service.NewTaskService(taskRepo)
	emailService := service.NewEmailService(&conf)
	karmaService := service.NewKarmaService(userBadgeRepo, userRepo)
	interestScoreService := service.NewInterestScoreService(interestScoreRepo)
	tagPreferenceService := service.NewTagPreferenceService(tagPreferenceRepo)

	// Kafka
	producer := kafka.NewProducer(conf.Kafka)
	consumer := kafka.NewConsumer(conf.Kafka, emailService, karmaService, interestScoreService)

	// Handler
	taskHandler := handler.NewTaskHandler(taskService, producer)

	logger.Infof("[Worker] Kafka worker is running")

	// Run Producer in background - Process due tasks
	go func() {
		ctx := context.Background()
		for {
			taskHandler.ProcessDueTasks(ctx)
			time.Sleep(5 * time.Second)
		}
	}()

	// Run Tag Preference Updater in background (daily at 2 AM UTC)
	go func() {
		for {
			// Calculate time until next 2 AM UTC
			now := time.Now().UTC()
			next2AM := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, time.UTC)
			if now.After(next2AM) {
				next2AM = next2AM.Add(24 * time.Hour)
			}
			timeUntil2AM := next2AM.Sub(now)

			logger.Infof("[TagPreference] Next update scheduled at %s (in %s)", next2AM.Format("2006-01-02 15:04:05"), timeUntil2AM)

			time.Sleep(timeUntil2AM)

			logger.Infof("[TagPreference] Starting daily tag preference update...")
			if err := tagPreferenceService.UpdateAllActiveUsers(); err != nil {
				logger.Errorf("[Error] Tag preference update failed: %v", err)
			} else {
				logger.Infof("[TagPreference] Daily update completed successfully")
			}
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
		logger.Errorf("[ERROR] Close postgresql fail: %s", err)
	}
}
