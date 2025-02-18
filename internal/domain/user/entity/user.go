package entity

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	Role      int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
