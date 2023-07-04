package caches

import (
	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(dns string) *redis.Client {
	opt, err := redis.ParseURL("redis://@localhost:6379/0")
	if err != nil {
		logger.Panicf("Error parsing redis url: %v", err)
	}

	return redis.NewClient(opt)
}
