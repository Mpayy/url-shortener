package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisClient struct {
	Client *redis.Client
	ctx    context.Context
}

func NewRedis(config *viper.Viper) *RedisClient {
	addr := fmt.Sprintf("%s:%d", config.GetString("REDIS_HOST"), config.GetInt("REDIS_PORT"))
	db := config.GetInt("REDIS_DB")

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	return &RedisClient{
		Client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
