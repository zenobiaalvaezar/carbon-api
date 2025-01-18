package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres() {
	dsn := os.Getenv("DATABASE_URL")

	// dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require&PreferSimpleProtocol=true&statement_cache_mode=describe", dbuser, dbpass, dbhost, dbport, dbname)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	// set max idle connection
	sqlDB, err := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
}

func ClosePostgres() {
	if DB != nil {
		db, err := DB.DB()
		if err != nil {
			log.Printf("Error getting DB instance: %v", err)
			return
		}

		if err := db.Close(); err != nil {
			log.Printf("Error closing DB: %v", err)
			return
		}

		log.Println("Connection to database closed")
	}
}
