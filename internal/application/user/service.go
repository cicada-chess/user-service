package user

import (
	"context"
	"errors"

	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrUsernameExists = errors.New("username already exists")
	ErrUserNotFound   = errors.New("user not found")
)

type userService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) interfaces.UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if dbUser, err := u.repo.GetByEmail(ctx, user.Email); err != nil {
		return nil, err
	} else if dbUser != nil {
		return nil, ErrEmailExists
	}

	if dbUser, err := u.repo.GetByUsername(ctx, user.Username); err != nil {
		return nil, err
	} else if dbUser != nil {
		return nil, ErrUsernameExists
	}

	if err := entity.ValidatePassword(user.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := entity.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	createdUser, err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (u *userService) GetById(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *userService) UpdateInfo(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, err := u.repo.GetById(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if err := entity.ValidatePassword(user.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := entity.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	updatedUser, err := u.repo.UpdateInfo(ctx, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil

}

func (u *userService) Delete(ctx context.Context, id string) error {
	if dbUser, err := u.repo.GetById(ctx, id); err != nil {
		return err
	} else if dbUser == nil {
		return ErrUserNotFound
	}

	err := u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil

}

func (u *userService) GetAll(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}

func (u *userService) ChangePassword(ctx context.Context, id, old_password, new_password string) error {
	return nil
}

func (u *userService) ToggleActive(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func (u *userService) GetRating(ctx context.Context, id string) (int, error) {
	return 0, nil
}

func (u *userService) UpdateRating(ctx context.Context, id string, delta int) (int, error) {
	return 0, nil
}
