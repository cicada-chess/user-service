package user

import (
	"context"
	"errors"

	"github.com/lib/pq"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
)

var (
	ErrEmailExists       = errors.New("email already exists")
	ErrUsernameExists    = errors.New("username already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidUUIDFormat = errors.New("invalid UUID format")
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
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "22P02" {
			return nil, ErrInvalidUUIDFormat
		}
		return nil, err
	} else if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *userService) UpdateInfo(ctx context.Context, user *entity.User) (*entity.User, error) {
	exists, err := u.repo.CheckUserExists(ctx, user.ID)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, ErrUserNotFound
	}

	if user.Password != "" {
		if err := entity.ValidatePassword(user.Password); err != nil {
			return nil, err
		}

		hashedPassword, err := entity.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	} else {
		user.Password, _ = u.repo.GetPasswordById(ctx, user.ID)
	}
	updatedUser, err := u.repo.UpdateInfo(ctx, user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil

}

func (u *userService) Delete(ctx context.Context, id string) error {
	exists, err := u.repo.CheckUserExists(ctx, id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "22P02" {
			return ErrInvalidUUIDFormat
		}
		return err
	} else if !exists {
		return ErrUserNotFound
	}

	err = u.repo.Delete(ctx, id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "22P02" {
			return ErrInvalidUUIDFormat
		}
		return err
	}

	return nil

}

func (u *userService) GetAll(ctx context.Context, page, limit, search, sort_by, order string) ([]*entity.User, error) {
	users, err := u.repo.GetAll(ctx, page, limit, search, sort_by, order)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userService) ChangePassword(ctx context.Context, id, old_password, new_password string) error {
	exists, err := u.repo.CheckUserExists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return ErrUserNotFound
	}

	if err := entity.ValidatePassword(new_password); err != nil {
		return err
	}
	dbPassword, _ := u.repo.GetPasswordById(ctx, id)

	if valid := entity.ComparePasswords(dbPassword, old_password); !valid {
		return ErrInvalidPassword
	}

	hashedPassword, err := entity.HashPassword(new_password)
	if err != nil {
		return err
	}

	err = u.repo.ChangePassword(ctx, id, hashedPassword)

	if err != nil {
		return err
	}

	return nil
}

func (u *userService) ToggleActive(ctx context.Context, id string) (bool, error) {
	exists, err := u.repo.CheckUserExists(ctx, id)
	if err != nil {
		return false, err
	} else if !exists {
		return false, ErrUserNotFound
	}

	isActive, err := u.repo.ToggleActive(ctx, id)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func (u *userService) GetRating(ctx context.Context, id string) (int, error) {
	exists, err := u.repo.CheckUserExists(ctx, id)
	if err != nil {
		return 0, err
	} else if !exists {
		return 0, ErrUserNotFound
	}

	rating, err := u.repo.GetRating(ctx, id)
	if err != nil {
		return 0, err
	}
	return rating, nil
}

func (u *userService) UpdateRating(ctx context.Context, id string, delta int) (int, error) {
	exists, err := u.repo.CheckUserExists(ctx, id)
	if err != nil {
		return 0, err
	} else if !exists {
		return 0, ErrUserNotFound
	}

	rating, err := u.repo.UpdateRating(ctx, id, delta)
	if err != nil {
		return 0, err
	}
	return rating, nil
}

func (u *userService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, ErrUserNotFound
	}
	user.Password, _ = u.repo.GetPasswordById(ctx, user.ID)
	return user, nil
}

func (u *userService) UpdatePasswordById(ctx context.Context, id, password string) error {
	if err := entity.ValidatePassword(password); err != nil {
		return err
	}

	hashedPassword, err := entity.HashPassword(password)
	if err != nil {
		return err
	}

	err = u.repo.ChangePassword(ctx, id, hashedPassword)

	if err != nil {
		return err
	}

	return nil
}
