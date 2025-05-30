package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"shoeshop/product-service/internal/model"
)

type Cache interface {
	Get(ctx context.Context, key string) (*model.Product, error)
	Set(ctx context.Context, key string, product *model.Product) error
	Delete(ctx context.Context, key string) error
	Close() error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Проверяем подключение
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &redisCache{
		client: client,
	}, nil
}

func (c *redisCache) Get(ctx context.Context, key string) (*model.Product, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Кэш пуст
		}
		return nil, err
	}

	var product model.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *redisCache) Set(ctx context.Context, key string, product *model.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	// Кэшируем на 1 час
	return c.client.Set(ctx, key, data, time.Hour).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *redisCache) Close() error {
	return c.client.Close()
} 