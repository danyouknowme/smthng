package datasources

import (
	"github.com/danyouknowme/smthng/internal/config"
	"github.com/danyouknowme/smthng/internal/datasources/caches"
	"github.com/danyouknowme/smthng/internal/datasources/drivers"
	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataSources interface {
	GetRedisClient() *redis.Client
}

type datasources struct {
	redis *redis.Client
	mongo *mongo.Client
}

func NewDataSources(config *config.AppConfig) DataSources {
	redisClient, err := caches.NewRedisClient(config.RedisURI)
	if err != nil {
		logger.Panicf("Failed to connect to redis: %v", err)
	}

	mongoClient, err := drivers.NewMongoClient(config.MongoURI)
	if err != nil {
		logger.Panicf("Failed to connect to mongo: %v", err)
	}

	return &datasources{
		redis: redisClient,
		mongo: mongoClient,
	}
}

func (ds *datasources) GetRedisClient() *redis.Client {
	return ds.redis
}
