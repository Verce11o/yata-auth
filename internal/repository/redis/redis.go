package redis

import (
	"github.com/Verce11o/yata-auth/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) *redis.Client {
	redisHost := cfg.Redis.Host

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return client
}
