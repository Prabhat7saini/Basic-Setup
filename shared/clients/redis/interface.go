// Package redis - client.go
package redis

import (
	"context"
	
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}) error 
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error // Added Close method for proper cleanup
	SetWithExp(ctx context.Context, key string, value interface{}, minutes int) error



	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)

	HSet(ctx context.Context, key string, values map[string]interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HDel(ctx context.Context, key string, fields ...string) error
}
