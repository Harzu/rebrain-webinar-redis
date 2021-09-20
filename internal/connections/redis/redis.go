package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/Harzu/rebrain-webinar-redis/internal/config"
)

const (
	redisPrefix       = "rebrain_webinar"
	readRedisTimeOut  = 200 * time.Millisecond
	writeRedisTimeOut = 200 * time.Millisecond
)

type Client interface {
	GetConnection() *redis.Client
	BuildKeyWithPrefix(key string) string
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

func (r *client) BuildKeyWithPrefix(key string) string {
	return fmt.Sprintf("%s_%s", redisPrefix, key)
}
