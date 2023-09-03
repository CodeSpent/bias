package config

import (
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
)

var RedisClient *redis.Client

func Initialize() error {
	godotenv.Load()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_SERVER_ADDR"),
	})

	godotenv.Load()

	return nil
}

func Close() {
	if RedisClient != nil {
		_ = RedisClient.Close()
	}
}
