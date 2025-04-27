package interfaces

import (
	"context"

	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error)
	GetByUserID(ctx context.Context, userID string) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error)
}
