package usecases

import (
	"context"
	"time"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type channelUsecase struct {
	channelRepository repositories.ChannelRepository
}

type ChannelUsecase interface {
	GetChannelByID(ctx context.Context, ID string) (*domains.Channel, error)
	CreateNewChannel(ctx context.Context, channel *domains.CreateChannelRequest) error
	IsMember(ctx context.Context, channelID string, userID string) bool
}

func NewChannelUsecase(channelRepository repositories.ChannelRepository) ChannelUsecase {
	return &channelUsecase{
		channelRepository: channelRepository,
	}
}

func (usecase *channelUsecase) GetChannelByID(ctx context.Context, ID string) (*domains.Channel, error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	channel, err := usecase.channelRepository.FindById(ctx, objID)
	if err != nil {
		return nil, err
	}

	return channel.Serialize(), nil
}

func (usecase *channelUsecase) CreateNewChannel(ctx context.Context, channel *domains.CreateChannelRequest) error {
	members := []primitive.ObjectID{}
	for _, memberID := range channel.Members {
		objID, err := primitive.ObjectIDFromHex(memberID)
		if err != nil {
			return err
		}
		members = append(members, objID)
	}

	var guildID *string
	if channel.GuildID == "" {
		guildID = nil
	} else {
		guildID = &channel.GuildID
	}

	var isPublic bool
	if channel.IsDM {
		isPublic = false
	} else {
		isPublic = channel.IsPublic
	}

	channelMongo := &domains.ChannelMongo{
		GuildID:      guildID,
		Name:         channel.Name,
		IsPublic:     isPublic,
		IsDM:         channel.IsDM,
		LastActivity: time.Now(),
		Members:      members,
		Messages:     []primitive.ObjectID{},
	}
	return usecase.channelRepository.Create(ctx, channelMongo)
}

func (usecase *channelUsecase) IsMember(ctx context.Context, channelID string, userID string) bool {
	return true
}
