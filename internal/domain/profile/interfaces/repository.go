package interfaces

import (
	"context"

	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
)

// ProfileRepository определяет методы для работы с профилями пользователей
type ProfileRepository interface {
	Create(ctx context.Context, profile *entity.Profile) (*entity.Profile, error)
	GetByUserID(ctx context.Context, userID string) (*entity.Profile, error)
	Update(ctx context.Context, profile *entity.Profile) (*entity.Profile, error)
	SaveAvatar(ctx context.Context, userID string, avatarData []byte, fileType string) (string, error)
}
