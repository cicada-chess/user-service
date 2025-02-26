package interfaces

import (
	"context"

	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetById(ctx context.Context, id string) (*entity.User, error)
	UpdateInfo(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*entity.User, error)
	ChangePassword(ctx context.Context, id, old_password, new_password string) error
	ToggleActive(ctx context.Context, id string) (bool, error)
	GetRating(ctx context.Context, id string) (int, error)
	UpdateRating(ctx context.Context, id string, delta int) (int, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}
