package usecases

import (
	"context"
	"time"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messageUsecase struct {
	messageRepository repositories.MessageRepository
	userRepository    repositories.UserRepository
}

type MessageUsecase interface {
	CreateNewMessage(ctx context.Context, message *domains.CreateMessageRequest) (*domains.Message, error)
}

func NewMessageUsecase(messageRepository repositories.MessageRepository, userRepository repositories.UserRepository) MessageUsecase {
	return &messageUsecase{
		messageRepository: messageRepository,
		userRepository:    userRepository,
	}
}

func (usecase *messageUsecase) CreateNewMessage(ctx context.Context, message *domains.CreateMessageRequest) (*domains.Message, error) {
	userObjectID, err := primitive.ObjectIDFromHex(message.UserID)
	if err != nil {
		return nil, err
	}

	channelObjectID, err := primitive.ObjectIDFromHex(message.ChannelID)
	if err != nil {
		return nil, err
	}

	user, err := usecase.userRepository.FindByID(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	messageMongo := &domains.MessageMongo{
		ID:        primitive.NewObjectID(),
		Text:      message.Text,
		UserID:    user.ID,
		ChannelID: channelObjectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := usecase.messageRepository.Create(ctx, messageMongo); err != nil {
		return nil, err
	}

	newMessage := messageMongo.Serialize()
	newMessage.Member = *user.Serialize()

	return newMessage, nil
}
