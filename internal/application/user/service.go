package user

import (
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

func (u *userService) GetById(id string) (*entity.User, error) {
	user, err := u.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) Create(user *entity.User) (string, error) {
	id, err := u.repo.Create(user)
	if err != nil {
		return "", err
	}
	return id, nil
}
