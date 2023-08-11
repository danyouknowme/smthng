package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Tag          string `json:"tag"`
	ProfileImage string `json:"profile_image"`
	IsOnline     bool   `json:"is_online"`
	Friends      []User `json:"friends"`
	Request      []User `json:"requests"`
}

type UserMongo struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Username     string               `bson:"username"`
	Email        string               `bson:"email"`
	Password     string               `bson:"password"`
	Tag          string               `bson:"tag"`
	ProfileImage string               `bson:"profile_image"`
	IsOnline     bool                 `bson:"is_online"`
	Friends      []primitive.ObjectID `bson:"friends"`
	Requests     []primitive.ObjectID `bson:"requests"`
	CreatedAt    time.Time            `bson:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *UserMongo) Serialize() *User {
	return &User{
		ID:           user.ID.Hex(),
		Username:     user.Username,
		Email:        user.Email,
		Tag:          user.Tag,
		ProfileImage: user.ProfileImage,
		IsOnline:     user.IsOnline,
	}
}

func (user *UserMongo) SerializeToMember(isFriend bool) *Member {
	return &Member{
		ID:           user.ID.Hex(),
		Username:     user.Username,
		ProfileImage: user.ProfileImage,
		IsOnline:     user.IsOnline,
		IsFriend:     isFriend,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
