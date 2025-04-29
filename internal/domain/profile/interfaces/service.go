package interfaces

import (
	"context"
	"mime/multipart"

	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
)

type ProfileService interface {
	CreateProfile(ctx context.Context, userID string) (*entity.Profile, error)
	GetProfile(ctx context.Context, userID string) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error)
	UploadAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (string, error)
	GetUserIDFromToken(ctx context.Context, token string) (string, error)
}
