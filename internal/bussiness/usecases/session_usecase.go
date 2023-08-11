package usecases

import (
	"context"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type sessionUsecase struct {
	sessionRepository repositories.SessionRepository
}

type SessionUsecase interface {
	CreateNewSession(ctx context.Context, session *domains.SessionMongo) (*domains.Session, error)
}

func NewSessionUsecase(sessionRepository repositories.SessionRepository) SessionUsecase {
	return &sessionUsecase{
		sessionRepository: sessionRepository,
	}
}

func (usecase *sessionUsecase) CreateNewSession(ctx context.Context, session *domains.SessionMongo) (*domains.Session, error) {
	session.ID = primitive.NewObjectID()

	if err := usecase.sessionRepository.Create(ctx, session); err != nil {
		return nil, err
	}

	return session.Serialize(), nil
}
