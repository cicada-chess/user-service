package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (string, error) {
	return "", nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*entity.User, error) {
	return nil, nil
}

func (r *userRepository) UpdateInfo(ctx context.Context, user *entity.User) (*entity.User, error) {
	return nil, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	return nil, nil
}

func (r *userRepository) ChangePassword(ctx context.Context, id, old_password, new_password string) error {
	return nil
}

func (r *userRepository) ToggleActive(ctx context.Context, id string) (bool, error) {
	return false, nil
}

func (r *userRepository) GetRating(ctx context.Context, id string) (int, error) {
	return 0, nil
}

func (r *userRepository) UpdateRating(ctx context.Context, id string, delta int) (int, error) {

	return 0, nil
}
