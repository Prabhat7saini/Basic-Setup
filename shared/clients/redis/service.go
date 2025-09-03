package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	rdb *redis.Client
}

func (r *redisClient) Set(ctx context.Context, key string, value interface{}) error {
	return r.rdb.Set(ctx, key, value, 0).Err()
}

func (r *redisClient) SetWithExp(ctx context.Context, key string, value interface{}, minutes int) error {
	exp := time.Duration(minutes) * time.Minute
	return r.rdb.Set(ctx, key, value, exp).Err()
}
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *redisClient) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *redisClient) Exists(ctx context.Context, key string) (bool, error) {
	val, err := r.rdb.Exists(ctx, key).Result()
	return val > 0, err
}
func (r *redisClient) Close() error {
	return r.rdb.Close()
}


func (r *redisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.rdb.SAdd(ctx, key, members...).Err()
}

func (r *redisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.rdb.SRem(ctx, key, members...).Err()
}

func (r *redisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.rdb.SMembers(ctx, key).Result()
}

// --- Hashes ---
func (r *redisClient) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	return r.rdb.HSet(ctx, key, values).Err()
}

func (r *redisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return r.rdb.HGet(ctx, key, field).Result()
}

func (r *redisClient) HDel(ctx context.Context, key string, fields ...string) error {
	return r.rdb.HDel(ctx, key, fields...).Err()
}