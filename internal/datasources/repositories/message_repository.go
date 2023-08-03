package repositories

import (
	"context"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageRepository struct {
	collection *mongo.Collection
}

type MessageRepository interface {
	Create(ctx context.Context, message *domains.MessageMongo) error
}

func NewMessageRepository(ds datasources.DataSources) MessageRepository {
	return &messageRepository{
		collection: ds.GetMongoCollection("messages"),
	}
}

func (r *messageRepository) Create(ctx context.Context, message *domains.MessageMongo) error {
	_, err := r.collection.InsertOne(ctx, message)
	return err
}
