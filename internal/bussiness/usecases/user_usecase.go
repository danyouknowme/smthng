package usecases

import (
	"context"
	"time"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources/repositories"
	"github.com/danyouknowme/smthng/pkg/apperrors"
	"github.com/danyouknowme/smthng/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userRepository repositories.UserRepository
}

type UserUsecase interface {
	CreateNewUser(ctx context.Context, user *domains.RegisterRequest) error
	Authenticate(ctx context.Context, req *domains.LoginRequest) (string, error)
}

func NewUserUsecase(userRepository repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (usecase *userUsecase) CreateNewUser(ctx context.Context, user *domains.RegisterRequest) error {
	userMongo := &domains.UserMongo{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Tag:       helpers.GenerateTag(),
		Friends:   []primitive.ObjectID{},
		Requests:  []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for {
		_, err := usecase.userRepository.FindByTag(ctx, userMongo.Tag)
		if err != nil {
			err := usecase.userRepository.Create(ctx, userMongo)
			return err
		}

		userMongo.Tag = helpers.GenerateTag()
	}
}

func (usecase *userUsecase) Authenticate(ctx context.Context, req *domains.LoginRequest) (string, error) {
	userMongo, err := usecase.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return "", err
	}

	if err = helpers.CheckPassword(req.Password, userMongo.Password); err != nil {
		return "", apperrors.ErrInvalidUsernameOrPassword
	}

	return userMongo.ID.Hex(), nil
}
