package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Prabhat7saini/Basic-Setup/config"
	"github.com/redis/go-redis/v9"
)

var (
	instance Client
	once     sync.Once
)

func InitRedis(cfg *config.Env) (Client, error) {
	var err error

	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Db,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if pingErr := rdb.Ping(ctx).Err(); pingErr != nil {
			err = fmt.Errorf("failed to connect to Redis at %s: %w", cfg.Redis.Addr, pingErr)
			return
		}

		instance = &redisClient{rdb: rdb}
	})

	if instance == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("redis client is nil (unknown error)")
	}
	fmt.Println("redis client is initialized")
	return instance, nil
}

func GetRedisClient() (Client, error) {
	if instance == nil {
		return nil, fmt.Errorf("redis client not initialized ")
	}
	return instance, nil
}
func CloseRedis() error {
	if instance == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// Cast to *redisClient (your wrapper struct)
	rc, ok := instance.(*redisClient)
	if !ok {
		return fmt.Errorf("invalid redis client type")
	}

	if err := rc.rdb.Close(); err != nil {
		return fmt.Errorf("failed to close redis client: %w", err)
	}

	instance = nil
	return nil
}
