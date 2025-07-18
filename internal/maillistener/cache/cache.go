package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("key not found")
	ErrExpired  = errors.New("key expired")
)

type Cacher interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Has(ctx context.Context, key string) (bool, error)
}
