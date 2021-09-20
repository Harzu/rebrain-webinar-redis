package locker

import (
	"context"
	"time"

	"github.com/bsm/redislock"

	"github.com/Harzu/rebrain-webinar-redis/internal/connections/redis"
)

type Locker interface {
	Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error)
	ObtainLinear(ctx context.Context, key string, ttl, interval time.Duration) (Lock, error)
}

type locker struct {
	redisClient redis.Client
}

func New(redisClient redis.Client) Locker {
	return &locker{redisClient: redisClient}
}

func (l *locker) Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	redisLocker := redislock.New(l.redisClient.GetConnection())
	redisLock, err := redisLocker.Obtain(ctx, l.redisClient.BuildKeyWithPrefix(key), ttl, nil)
	return &lock{redisLock: redisLock}, err
}

func (l *locker) ObtainLinear(ctx context.Context, key string, ttl, interval time.Duration) (Lock, error) {
	redisLocker := redislock.New(l.redisClient.GetConnection())
	redisLock, err := redisLocker.Obtain(ctx, l.redisClient.BuildKeyWithPrefix(key), ttl, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(interval),
	})
	return &lock{redisLock: redisLock}, err
}
