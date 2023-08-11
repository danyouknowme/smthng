package domains

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/danyouknowme/smthng/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Text      *string            `bson:"text,omitempty"`
	File      *string            `bson:"file,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ChannelID primitive.ObjectID `bson:"channel_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Message struct {
	ID        string    `json:"id"`
	Text      *string   `json:"text"`
	File      *string   `json:"file"`
	ChannelID string    `json:"channel_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Member    *Member   `json:"member"`
}

func (message *MessageMongo) Serialize() *Message {
	return &Message{
		ID:        message.ID.Hex(),
		Text:      message.Text,
		File:      message.File,
		ChannelID: message.ChannelID.Hex(),
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

type CreateMessageRequest struct {
	Text      *string               `json:"text"`
	File      *multipart.FileHeader `json:"file"`
	UserID    string                `json:"user_id"`
	ChannelID string                `json:"channel_id"`
}

type MessageRequest struct {
	Text *string               `form:"text"`
	File *multipart.FileHeader `form:"file"`
}

func (message *MessageRequest) Validate() error {
	if message.Text == nil && message.File == nil {
		return fmt.Errorf("text or file must be provided")
	}

	if message.Text != nil {
		if len(strings.TrimSpace(*message.Text)) == 0 {
			return fmt.Errorf("text cannot be empty")
		}
		if len(*message.Text) > 1000 {
			return fmt.Errorf("text too long")
		}
	}

	if message.File != nil {
		out, err := message.File.Open()
		if err != nil {
			return err
		}
		defer out.Close()

		mineType := helpers.GetFileContentType(out)

		if !strings.Contains(mineType, "image") {
			return fmt.Errorf("file type not supported")
		}

		if message.File.Size > 5e+6 {
			return fmt.Errorf("file size too large")
		}
	}

	return nil
}
