package config

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisUrl := os.Getenv("REDIS_URL")

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	log.Println("Connected to Redis")
}

func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}
}
