package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionMongo struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       string             `bson:"user_id"`
	RefreshToken string             `bson:"refresh_token"`
	UserAgent    string             `bson:"user_agent"`
	ClientIP     string             `bson:"client_ip"`
	ExpiredAt    time.Time          `bson:"expired_at"`
}

type Session struct {
	ID           string    `json:"session_id"`
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	ExpiredAt    time.Time `json:"expired_at"`
}

func (session *SessionMongo) Serialize() *Session {
	return &Session{
		ID:           session.ID.Hex(),
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIP,
		ExpiredAt:    session.ExpiredAt,
	}
}
