package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file is found")
	}
}

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database %v", err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.Db %v", err)
	}
	sqlDb.SetMaxOpenConns(25)
	sqlDb.SetMaxIdleConns(25)
	sqlDb.SetConnMaxLifetime(5 * time.Minute)
	fmt.Println("Database successfully connected to GORM")
	return db, nil
}
