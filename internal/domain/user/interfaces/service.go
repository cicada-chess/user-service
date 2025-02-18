package interfaces

import (
	"context"

	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
)

type UserService interface {
	Create(ctx context.Context, user *entity.User) (string, error)
	GetById(ctx context.Context, id string) (*entity.User, error)
}
