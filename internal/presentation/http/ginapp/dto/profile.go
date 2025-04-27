package dto

import "time"

type Profile struct {
	UserID      string    `json:"user_id"`
	Description string    `json:"description"`
	Age         int       `json:"age"`
	Location    string    `json:"location"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Description *string `json:"description"`
	Age         *int    `json:"age"`
	Location    *string `json:"location"`
}
