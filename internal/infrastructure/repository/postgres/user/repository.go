package user

import (
	"context"
	"database/sql"

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

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	createdUser, err := r.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	createdUser.Password = ""
	return createdUser, nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	user := &entity.User{}
	err := r.db.Get(user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
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

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	user := &entity.User{}
	err := r.db.Get(user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	user := &entity.User{}
	err := r.db.Get(user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
