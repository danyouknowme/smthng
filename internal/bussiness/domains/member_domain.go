package domains

import "time"

type MemberMongo struct{}

type Member struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	ProfileImage string    `json:"profile_image"`
	IsOnline     bool      `json:"is_online"`
	Nickname     *string   `json:"nickname"`
	Color        *string   `json:"color"`
	IsFriend     bool      `json:"is_friend"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
