package userstorage

import (
	"context"
	"fmt"

	"github.com/go-redis/cache/v8"

	"github.com/Harzu/rebrain-webinar-redis/internal/entities"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/connections/redis"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/constants"
	"github.com/Harzu/rebrain-webinar-redis/internal/system/localcache/lru"
)

const storageKey = "user"

type Storage interface {
	Set(ctx context.Context, user *entities.User) error
	Get(ctx context.Context, userID int64) (*entities.User, error)
	Del(ctx context.Context, userID int64) error
}

type storage struct {
	cache *cache.Cache
}

func New(redisClient redis.Client, lruSize int) (Storage, error) {
	lruCache, err := lru.New(lruSize)
	if err != nil {
		return nil, err
	}

	cacheStorage := cache.New(&cache.Options{
		Redis:      redisClient.GetConnection(),
		LocalCache: lruCache,
	})

	return &storage{cache: cacheStorage}, nil
}

func (s *storage) Set(ctx context.Context, user *entities.User) error {
	return s.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   s.buildKey(user.ID),
		Value: user,
	})
}

func (s *storage) Get(ctx context.Context, userID int64) (user *entities.User, _ error) {
	if err := s.cache.Get(ctx, s.buildKey(userID), user); err != nil {
		return nil, fmt.Errorf("failed to get user from cache")
	}
	return user, nil
}

func (s *storage) Del(ctx context.Context, userID int64) error {
	return s.cache.Delete(ctx, s.buildKey(userID))
}

func (s *storage) buildKey(userID int64) string {
	return fmt.Sprintf("%s_%s_%d", constants.RedisPrefix, storageKey, userID)
}
