package interfaces

import (
	"context"
	"io"
)

type ProfileStorage interface {
	SaveAvatar(ctx context.Context, userID string, avatarData io.Reader, fileExt string) (string, error)
}
