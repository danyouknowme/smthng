package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChannelMongo struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	GuildID      *string              `bson:"guild_id,omitempty"`
	Name         string               `bson:"name"`
	IsPublic     bool                 `bson:"is_public"`
	IsDM         bool                 `bson:"is_dm"`
	LastActivity time.Time            `bson:"last_activity"`
	Members      []primitive.ObjectID `bson:"members"`
	Messages     []primitive.ObjectID `bson:"messages"`
}

type Channel struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	IsPublic        bool      `json:"is_public"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	HasNotification bool      `json:"has_notification"`
}

func (channel *ChannelMongo) Serialize() *Channel {
	return &Channel{
		ID:              channel.ID.Hex(),
		Name:            channel.Name,
		IsPublic:        channel.IsPublic,
		CreatedAt:       channel.LastActivity,
		UpdatedAt:       channel.LastActivity,
		HasNotification: false,
	}
}

type CreateChannelRequest struct {
	GuildID  string   `json:"guild_id"`
	Name     string   `json:"name" binding:"required"`
	IsPublic bool     `json:"is_public"`
	IsDM     bool     `json:"is_dm"`
	Members  []string `json:"members"`
}
