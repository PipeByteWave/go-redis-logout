package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var Redis *RedisClient

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func InitRedis(addr string) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error conectando a Redis:", err)
	}

	Redis = &RedisClient{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisClient) SetValue(key, value string) error {
	return r.client.Set(r.ctx, key, value, 24*time.Hour).Err()
}

func (r *RedisClient) GetValue(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) IsTokenBlacklisted(token string) (bool, error) {
	result, err := r.GetValue("blacklist:" + token)
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error al verificar el token en Redis: %v", err)
	}
	return result == "true", nil
}
