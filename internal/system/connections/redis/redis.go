package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/Harzu/rebrain-webinar-redis/internal/system/config"
)

const (
	readRedisTimeOut  = 200 * time.Millisecond
	writeRedisTimeOut = 200 * time.Millisecond
)

type Client interface {
	GetConnection() *redis.Client
}

type client struct {
	client *redis.Client
}

func New(cfg *config.Redis) (*client, error) {
	redisClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: []string{cfg.URL},
		ReadTimeout:   readRedisTimeOut,
		WriteTimeout:  writeRedisTimeOut,
		Password:      cfg.Password,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &client{client: redisClient}, nil
}

func (r *client) GetConnection() *redis.Client {
	return r.client
}

func (r *client) Close() error {
	return r.client.Close()
}
