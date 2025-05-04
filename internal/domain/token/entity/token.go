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
		"userId":     userId,
		"token_type": "account_confirmation",
		"expires_at": time.Now().Add(AccountConfirmationTTL).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func GeneratePasswordResetToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     userId,
		"token_type": "password_reset",
		"expires_at": time.Now().Add(AccountConfirmationTTL).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func ValidateToken(tokenString string, tokenType TokenType) (*string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalidOrExpired
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalidOrExpired
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, ErrTokenInvalidOrExpired
	}

	expires_at, ok := claims["expires_at"].(float64)
	if !ok || expires_at < float64(time.Now().Unix()) {
		return nil, ErrTokenInvalidOrExpired
	}

	token_type, ok := claims["token_type"].(string)
	if !ok || token_type != string(tokenType) {
		return nil, ErrTokenInvalidOrExpired
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return nil, ErrTokenInvalidOrExpired
	}
	return &userId, nil
}
