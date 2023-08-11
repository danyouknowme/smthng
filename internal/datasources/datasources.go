package datasources

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/danyouknowme/smthng/internal/config"
	"github.com/danyouknowme/smthng/internal/datasources/caches"
	"github.com/danyouknowme/smthng/internal/datasources/drivers"
	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataSources interface {
	GetRedisClient() *redis.Client
	GetMongoClient() *mongo.Client
	GetCloudinaryClient() *cloudinary.Cloudinary
	GetMongoCollection(collection string) *mongo.Collection
	Close() error
}

type datasources struct {
	redis      *redis.Client
	mongo      *mongo.Client
	cloudinary *cloudinary.Cloudinary
}

const DB = "smthng"

func NewDataSources(config *config.AppConfig) DataSources {
	redisClient, err := caches.NewRedisClient(config.RedisURI)
	if err != nil {
		logger.Panicf("Failed to connect to redis: %v", err)
	}

	mongoClient, err := drivers.NewMongoClient(config.MongoURI)
	if err != nil {
		logger.Panicf("Failed to connect to mongo: %v", err)
	}

	cloudinary, err := drivers.NewCloudinaryClient(config.CloudinaryURI)
	if err != nil {
		logger.Panicf("Failed to create gridfs bucket: %v", err)
	}

	return &datasources{
		redis:      redisClient,
		mongo:      mongoClient,
		cloudinary: cloudinary,
	}
}

func (ds *datasources) GetRedisClient() *redis.Client {
	return ds.redis
}

func (ds *datasources) GetMongoClient() *mongo.Client {
	return ds.mongo
}

func (ds *datasources) GetCloudinaryClient() *cloudinary.Cloudinary {
	return ds.cloudinary
}

func (ds *datasources) GetMongoCollection(collection string) *mongo.Collection {
	return ds.mongo.Database(DB).Collection(collection)
}

func (ds *datasources) Close() error {
	if err := ds.redis.Close(); err != nil {
		return err
	}
	if err := ds.mongo.Disconnect(context.Background()); err != nil {
		return err
	}

	return nil
}
