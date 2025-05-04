package interfaces

import (
	"context"
)

type NotificationSender interface {
	SendAccountConfirmation(ctx context.Context, userId, email, token string) error
	SendPasswordReset(ctx context.Context, userId, email, token string) error
}
