package caches

import (
	"context"
	"time"

	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(host string) (*redis.Client, error) {
	logger.Info("Registering redis client...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatalf("Failed to ping redis: %v", err)
		return nil, err
	}

	logger.Info("Registering redis client completed")
	return client, nil
}
