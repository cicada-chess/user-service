package entity

import (
	"time"
)

type Profile struct {
	UserID      string
	Username    string
	Description string
	Age         int
	Location    string
	AvatarURL   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Profile) IsValidAge() bool {
	return p.Age > 0
}
