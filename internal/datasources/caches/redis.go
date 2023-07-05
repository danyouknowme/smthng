package caches

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(host string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})
}
