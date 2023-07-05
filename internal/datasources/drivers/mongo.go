package drivers

import (
	"context"
	"time"

	"github.com/danyouknowme/smthng/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func NewMongoClient(dsn string) (*mongo.Client, error) {
	logger.Info("Registering mongo client...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		logger.Fatalf("Failed to create mongo client: %v", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatalf("Failed to ping mongo: %v", err)
		return nil, err
	}

	logger.Info("Registering mongo client completed")
	return client, nil
}
