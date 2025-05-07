package entity

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccountConfirmationTTL = 30 * 60 * time.Second
)

type TokenType string

const (
	AccountConfirmation TokenType = "account_confirmation"
	PasswordReset       TokenType = "password_reset"
)

var ErrTokenInvalidOrExpired = errors.New("token invalid or expired")

func GenerateAccountConfirmationToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userId,
		"token_type": string(AccountConfirmation),
		"expires_at": time.Now().Add(AccountConfirmationTTL).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func GeneratePasswordResetToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userId,
		"token_type": string(PasswordReset),
		"expires_at": time.Now().Add(AccountConfirmationTTL).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}
