package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messageRepository struct {
	collection *mongo.Collection
}

type MessageRepository interface {
	FindByID(ctx context.Context, ID primitive.ObjectID) (*domains.MessageMongo, error)
	Create(ctx context.Context, message *domains.MessageMongo) error
	UpdateByID(ctx context.Context, ID primitive.ObjectID, updatedText string) (*domains.MessageMongo, error)
}

func NewMessageRepository(ds datasources.DataSources) MessageRepository {
	return &messageRepository{
		collection: ds.GetMongoCollection("messages"),
	}
}

func (r *messageRepository) FindByID(ctx context.Context, ID primitive.ObjectID) (*domains.MessageMongo, error) {
	var message domains.MessageMongo
	err := r.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *messageRepository) Create(ctx context.Context, message *domains.MessageMongo) error {
	_, err := r.collection.InsertOne(ctx, message)
	return err
}

func (r *messageRepository) UpdateByID(ctx context.Context, ID primitive.ObjectID, updatedText string) (*domains.MessageMongo, error) {
	filter := bson.M{"_id": ID}
	update := bson.M{"$set": bson.M{
		"text":       updatedText,
		"updated_at": primitive.NewDateTimeFromTime(time.Now()),
	}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	fmt.Println("filter: ", filter)

	var updatedMessage domains.MessageMongo
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedMessage)
	if err != nil {
		return nil, err
	}

	return &updatedMessage, nil
}
