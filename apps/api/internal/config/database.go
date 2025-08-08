package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Neel-shetty/clarity/internal/domain"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
		return nil, fmt.Errorf("failed to connect to the database %v", err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.Db %v", err)
	}
	sqlDb.SetMaxOpenConns(25)
	sqlDb.SetMaxIdleConns(25)
	sqlDb.SetConnMaxLifetime(5 * time.Minute)
	fmt.Println("Database successfully connected to GORM")
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	return db, nil
}

var RedisClient *redis.Client

func ConnectRedisDB() (*redis.Client, error) {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("invalid REDIS_DB environment variable: must be a valid integer (got %q)", os.Getenv("REDIS_DB"))
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
		Protocol: 2,
	})

	ctx := context.Background()
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to the Redis Database")
	} else {
		log.Println("Connected to Redis")
	}
	return RedisClient, nil
}
