package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserMongo struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Username     string               `json:"username" bson:"username"`
	Email        string               `json:"email" bson:"email"`
	Password     string               `json:"password" bson:"password"`
	Tag          string               `json:"tag" bson:"tag"`
	ProfileImage string               `json:"profile_image" bson:"profile_image"`
	IsOnline     bool                 `json:"is_online" bson:"is_online"`
	Friends      []primitive.ObjectID `json:"friends" bson:"friends"`
	Requests     []primitive.ObjectID `json:"requests" bson:"requests"`
	CreatedAt    time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at" bson:"updated_at"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
