package entity

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	Role      int
	CreatedAt time.Time
	UpdatedAt time.Time
}
