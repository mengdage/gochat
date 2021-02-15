package storage

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"githum.com/mengdage/gochat/pkg/api"
)

type CacheStorage interface {
	SaveUserInfo(ctx context.Context, user *api.User) error
}

type cacheStorage struct {
	db *redis.Client
}

func SaveUserInfo(cache *redis.Client) CacheStorage {
	return &cacheStorage{
		db: cache,
	}
}

func (c *cacheStorage) SaveUserInfo(ctx context.Context, user *api.User) error {
	key := fmt.Sprintf("user:%d:session", user.ID)

	err := c.db.HSet(ctx, key, "id", user.ID, "name", user.Name).Err()
	if err != nil {
		return err
	}

	return nil
}
