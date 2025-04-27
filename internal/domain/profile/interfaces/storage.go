package interfaces

import "context"

type ProfileStorage interface {
	SaveAvatar(ctx context.Context, userID string, avatarData []byte, fileExt string) (string, error)
}
