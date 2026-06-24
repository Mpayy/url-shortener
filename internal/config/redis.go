package config

import (
	"context"
	"fmt"
	"time"
	"url-shortener/internal/exception"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const UrlCachePrefix = "url:"
const AuthPrefix = "auth:session:"

type RedisClient interface {
	Check(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RedisClientImpl struct {
	Client *redis.Client
}

func NewRedis(config *viper.Viper) RedisClient {
	addr := fmt.Sprintf("%s:%d", config.GetString("REDIS_HOST"), config.GetInt("REDIS_PORT"))
	db := config.GetInt("REDIS_DB")

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	return &RedisClientImpl{
		Client: rdb,
	}
}

func (r *RedisClientImpl) Check(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, exception.ErrInternalServer
	}

	if result == 0 {
		return false, exception.ErrUnauthorized
	}

	return true, nil
}

func (r *RedisClientImpl) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return exception.ErrInternalServer
	}
	return nil
}

func (r *RedisClientImpl) Delete(ctx context.Context, key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return exception.ErrInternalServer
	}
	return nil
}
