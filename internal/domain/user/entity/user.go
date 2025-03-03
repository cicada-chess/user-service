package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password is too short")
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

func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
