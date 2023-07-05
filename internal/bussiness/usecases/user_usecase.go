package usecases

import (
	"context"
	"time"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources/repositories"
	"github.com/danyouknowme/smthng/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userRepository repositories.UserRepository
}

type UserUsecase interface {
	CreateNewUser(ctx context.Context, user *domains.UserRequest) error
}

func NewUserUsecase(userRepository repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (usecase *userUsecase) CreateNewUser(ctx context.Context, user *domains.UserRequest) error {
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

	err := usecase.userRepository.Create(ctx, userMongo)
	return err
}
