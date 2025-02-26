package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	Role      int
	Rating    int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}
