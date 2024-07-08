package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Address  string
	Password string
	DBName   int
}

func NewRedisStorage(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DBName,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

type Storage struct {
	Client *redis.Client
}

func NewStorage(rdb *redis.Client) (*Storage, error) {
	return &Storage{Client: rdb}, nil
}

func (s *Storage) StoreUrl(url string) (bool, error) {
	ctx := context.Background()
	added, err := s.Client.SAdd(ctx, "image_urls0", url).Result()
	if err != nil {
		return false, fmt.Errorf("failed to add URL to Redis: %v", err)
	}
	return added == 1, nil
}
