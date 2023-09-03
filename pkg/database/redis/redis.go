package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClientService struct {
	client *redis.Client
}

func NewRedisClientService(addr, password string, db int) (*RedisClientService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClientService{
		client: client,
	}, nil
}

func (s *RedisClientService) GetClient() *redis.Client {
	return s.client
}
