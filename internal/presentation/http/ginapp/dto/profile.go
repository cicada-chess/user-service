package dto

import "time"

type Profile struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Age         int       `json:"age"`
	Location    string    `json:"location"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
