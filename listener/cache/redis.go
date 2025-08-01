package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, username, password string) (*RedisCache, error) {
	r := new(RedisCache)

	ropts := &redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
	}
	r.client = redis.NewClient(ropts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Ping(ctx).Err(); err != nil {
		r.client.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return r, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := c.client.Get(ctx, key).Bytes()

	if err == redis.Nil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return data, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()

	return err
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	return err
}

func (c *RedisCache) Has(ctx context.Context, key string) (bool, error) {
	_, err := c.client.Get(ctx, key).Bytes()

	if err == redis.Nil {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
