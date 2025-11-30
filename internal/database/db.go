package database

import (
	"fmt"
	"log"
	"os"
	"social-platform-kafka-worker/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var pgSingleton *gorm.DB

func InitPostgresql(conf *config.Config) {
	dbUrl := conf.Database.URL

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Fix: Sử dụng stdout thay vì nil
		logger.Config{
			SlowThreshold:             2 * time.Second,
			LogLevel:                  logger.Error, // Log only errors
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		panic("❌❌ Failed to connect database: " + err.Error())
	}
	fmt.Println("✅✅ Connect to the database successfully")
	pgSingleton = db
}

func GetDB() *gorm.DB {
	if pgSingleton == nil {
		panic("Connection to Database Postgres is not setup")
	}

	return pgSingleton
}

func ClosePostgresql() error {
	sqlDB, err := pgSingleton.DB()
	if err != nil {
		fmt.Println("failed to get sql.DB:", err)
	}
	defer sqlDB.Close()

	return err
}
