package repositories

import (
	"context"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type channelRepository struct {
	collection *mongo.Collection
}

type ChannelRepository interface {
	FindById(ctx context.Context, objID primitive.ObjectID) (*domains.ChannelMongo, error)
	Create(ctx context.Context, channel *domains.ChannelMongo) error
}

func NewChannelRepository(ds datasources.DataSources) ChannelRepository {
	return &channelRepository{
		collection: ds.GetMongoCollection("channels"),
	}
}

func (repo *channelRepository) FindById(ctx context.Context, objID primitive.ObjectID) (*domains.ChannelMongo, error) {
	var channel domains.ChannelMongo
	if err := repo.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (repo *channelRepository) Create(ctx context.Context, channel *domains.ChannelMongo) error {
	_, err := repo.collection.InsertOne(ctx, channel)
	if err != nil {
		return err
	}

	return nil
}
