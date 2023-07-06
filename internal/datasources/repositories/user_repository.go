package repositories

import (
	"context"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	collection *mongo.Collection
}

type UserRepository interface {
	Create(ctx context.Context, user *domains.UserMongo) error
	FindByTag(ctx context.Context, tag string) (*domains.UserMongo, error)
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	if _, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"tag": 1}, Options: options.Index().SetUnique(true)},
	}); err != nil {
		logger.Panicf("Error creating indexes in user collection: %v", err)
	}

	return &userRepository{
		collection: collection,
	}
}

func (repo *userRepository) Create(ctx context.Context, user *domains.UserMongo) error {
	_, err := repo.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) FindByTag(ctx context.Context, tag string) (*domains.UserMongo, error) {
	var user domains.UserMongo
	err := repo.collection.FindOne(ctx, bson.M{"tag": tag}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
