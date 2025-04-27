package dto

import "time"

type Profile struct {
	UserID      string    `json:"user_id" db:"user_id"`
	Description string    `json:"description" db:"description"`
	Age         int       `json:"age" db:"age"`
	Location    string    `json:"location" db:"location"`
	AvatarPath  string    `json:"avatar_url" db:"avatar_path"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
