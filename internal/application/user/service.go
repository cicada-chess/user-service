package user

import (
	"context"

	"github.com/google/uuid"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
)

type userService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) interfaces.UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) GetById(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) Create(ctx context.Context, user *entity.User) (string, error) {
	user.ID = uuid.New().String()
	user.Role = 1
	id, err := u.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

