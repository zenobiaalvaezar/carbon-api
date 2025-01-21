package config

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisUrl := os.Getenv("REDIS_URL")
	options, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	RedisClient = redis.NewClient(options)

	log.Println("Connected to Redis")
}

func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}
}
