package redis_dispatcher

import (
	"os"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	C *redis.Client
}

func NewRedisClient() *RedisClient {
	addr := os.Getenv("REDIS_HOST")
	if addr == "" {
		addr = "localhost:6379"
	}
	client := RedisClient{
		C: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
	return &client
}
