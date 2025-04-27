package entity

import (
	"time"
)

type Profile struct {
	UserID      string
	Description string
	Age         int
	Location    string
	AvatarPath  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Profile) IsValidAge() bool {
	return p.Age > 0
}
