package datasources

import (
	"github.com/danyouknowme/smthng/internal/config"
	"github.com/danyouknowme/smthng/internal/datasources/caches"
	"github.com/go-redis/redis/v8"
)

type DataSources interface {
	GetRedisClient() *redis.Client
}

type datasources struct {
	redis *redis.Client
}

func NewDataSources(config *config.AppConfig) DataSources {
	redisClient := caches.NewRedisClient(config.RedisHost)

	return &datasources{
		redis: redisClient,
	}
}

func (ds *datasources) GetRedisClient() *redis.Client {
	return ds.redis
}
