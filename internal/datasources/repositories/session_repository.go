package repositories

import (
	"context"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources"
	"go.mongodb.org/mongo-driver/mongo"
)

type sessionRepository struct {
	collection *mongo.Collection
}

type SessionRepository interface {
	Create(ctx context.Context, session *domains.SessionMongo) error
}

func NewSessionRepository(ds datasources.DataSources) SessionRepository {
	return &sessionRepository{
		collection: ds.GetMongoCollection("sessions"),
	}
}

func (repo *sessionRepository) Create(ctx context.Context, session *domains.SessionMongo) error {
	_, err := repo.collection.InsertOne(ctx, session)
	return err
}
