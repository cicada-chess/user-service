package interfaces

import "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"

type UserService interface {
	Create(user *entity.User) (string, error)
	GetById(id string) (*entity.User, error)
}
