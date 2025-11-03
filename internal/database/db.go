package database

import (
	"fmt"
	"social-platform-kafka-worker/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var pgSingleton *gorm.DB

func InitPostgresql(conf *config.Config) {
	dbUrl := conf.Database.URL

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	/* disable gorm logger
	db = db.Debug()
	*/

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
