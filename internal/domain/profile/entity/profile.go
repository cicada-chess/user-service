package entity

import (
	"time"
)

// Profile представляет дополнительную информацию о пользователе
type Profile struct {
	UserID      string    `json:"user_id"`
	Description string    `json:"description"`
	Age         int       `json:"age"`
	Location    string    `json:"location"`
	AvatarPath  string    `json:"avatar_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate проверяет правильность профиля
func (p *Profile) IsValidAge() bool {
	return p.Age > 0
}
