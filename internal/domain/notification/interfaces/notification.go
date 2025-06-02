package interfaces

import (
	"context"
)

type NotificationSender interface {
	SendAccountConfirmation(ctx context.Context, email, username, token string) error
	SendPasswordReset(ctx context.Context, email, username, token string) error
}
