package dto

import "time"

// Profile представляет DTO для профиля пользователя
type Profile struct {
	UserID      string    `json:"user_id"`
	Description string    `json:"description"`
	Age         int       `json:"age"`
	Location    string    `json:"location"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateProfileRequest представляет запрос на обновление профиля
type UpdateProfileRequest struct {
	Description string `json:"description"`
	Age         int    `json:"age"`
	Location    string `json:"location"`
}
