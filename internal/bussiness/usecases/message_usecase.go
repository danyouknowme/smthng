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
	GetMessageByID(ctx context.Context, messageID string) (*domains.Message, error)
	CreateNewMessage(ctx context.Context, message *domains.CreateMessageRequest) (*domains.Message, error)
	UpdateMessageByID(ctx context.Context, messageID, updatedText string) (*domains.Message, error)
}

func NewMessageUsecase(messageRepository repositories.MessageRepository, userRepository repositories.UserRepository) MessageUsecase {
	return &messageUsecase{
		messageRepository: messageRepository,
		userRepository:    userRepository,
	}
}

func (usecase *messageUsecase) GetMessageByID(ctx context.Context, messageID string) (*domains.Message, error) {
	messageObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, err
	}

	message, err := usecase.messageRepository.FindByID(ctx, messageObjectID)
	if err != nil {
		return nil, err
	}

	member, err := usecase.userRepository.FindByID(ctx, message.UserID)
	if err != nil {
		return nil, err
	}

	serializeMessage := message.Serialize()
	serializeMessage.Member = member.SerializeToMember(false)

	return serializeMessage, nil
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
	newMessage.Member = user.SerializeToMember(false)

	return newMessage, nil
}

func (usecase *messageUsecase) UpdateMessageByID(ctx context.Context, messageID, updatedText string) (*domains.Message, error) {
	messageObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, err
	}

	updatedMessage, err := usecase.messageRepository.UpdateByID(ctx, messageObjectID, updatedText)
	if err != nil {
		return nil, err
	}

	member, err := usecase.userRepository.FindByID(ctx, updatedMessage.UserID)
	if err != nil {
		return nil, err
	}

	serializeMessage := updatedMessage.Serialize()
	serializeMessage.Member = member.SerializeToMember(false)

	return serializeMessage, nil
}
