package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Text      string             `bson:"text"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ChannelID primitive.ObjectID `bson:"channel_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Message struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Member    User      `json:"user"`
}

func (message *MessageMongo) Serialize() *Message {
	return &Message{
		ID:        message.ID.Hex(),
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

type CreateMessageRequest struct {
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
}
