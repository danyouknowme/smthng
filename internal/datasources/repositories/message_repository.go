package repositories

import (
	"context"
	"fmt"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageRepository struct {
	collection *mongo.Collection
}

type MessageRepository interface {
	Create(ctx context.Context, message *domains.MessageMongo) error
	UpdateByID(ctx context.Context, ID primitive.ObjectID, updatedText string) (*domains.MessageMongo, error)
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

func (r *messageRepository) UpdateByID(ctx context.Context, ID primitive.ObjectID, updatedText string) (*domains.MessageMongo, error) {
	filter := bson.M{"_id": ID}
	update := bson.M{"$set": bson.M{"text": updatedText}}

	var updatedMessage domains.MessageMongo
	err := r.collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedMessage)
	if err != nil {
		return nil, err
	}

	res, _ := bson.MarshalExtJSON(updatedMessage, false, false)
	fmt.Println(string(res))

	return nil, nil
}
