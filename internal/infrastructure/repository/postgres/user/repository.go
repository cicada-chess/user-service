package user

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/dto"
)

const userFields = `id, username, email, role, rating, is_active, created_at, updated_at`

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `INSERT INTO users (username, email, password, rating, role, is_active) VALUES ($1, $2, $3, 0, 0, true)`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	createdUser, err := r.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT ` + userFields + ` FROM users WHERE id = $1`
	user := &dto.User{}
	err := r.db.Get(user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &entity.User{ID: user.ID, Username: user.Username, Email: user.Email, Rating: user.Rating, Role: user.Role, IsActive: user.IsActive, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (r *userRepository) UpdateInfo(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `UPDATE users SET username = $2, email = $3, password = $4, role = $5, rating = $6, is_active = $7, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.Role, user.Rating, user.IsActive)
	if err != nil {
		return nil, err
	}

	updatedUser, err := r.GetById(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetAll(ctx context.Context, page, limit, search, sort_by, order string) ([]*entity.User, error) {
	pageInt, err := strconv.Atoi(page)

	if err != nil {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt

	query := `SELECT ` + userFields + ` FROM users`

	if search != "" {
		query += ` WHERE email ILIKE $1 OR username ILIKE $1`
		search = "%" + search + "%"
	}

	if sort_by != "" {
		query += ` ORDER BY ` + sort_by
	} else {
		query += ` ORDER BY created_at`
	}

	if order != "" {
		if order != "asc" && order != "desc" {
			order = "asc"
		}
		query += ` ` + order
	} else {
		query += ` asc`
	}

	query += ` LIMIT $2 OFFSET $3`

	users := []*dto.User{}
	err = r.db.Select(&users, query, search, limitInt, offset)

	if err != nil {
		return nil, err
	}

	entityUsers := make([]*entity.User, len(users))
	for i, user := range users {
		entityUsers[i] = &entity.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Rating:    user.Rating,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	return entityUsers, nil
}

func (r *userRepository) ChangePassword(ctx context.Context, id, new_password string) error {
	query := `UPDATE users SET password = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id, new_password)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) ToggleActive(ctx context.Context, id string) (bool, error) {
	query := `UPDATE users SET is_active = NOT is_active, updated_at = NOW() WHERE id = $1 RETURNING is_active`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return false, err
	}

	updatedUser, err := r.GetById(ctx, id)
	if err != nil {
		return false, err
	}
	return updatedUser.IsActive, nil
}

func (r *userRepository) GetRating(ctx context.Context, id string) (int, error) {
	query := `SELECT rating FROM users WHERE id = $1`
	var rating int
	err := r.db.Get(&rating, query, id)
	if err != nil {
		return 0, err
	}
	return rating, nil
}

func (r *userRepository) UpdateRating(ctx context.Context, id string, delta int) (int, error) {
	query := `UPDATE users SET rating = rating + $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id, delta)
	if err != nil {
		return 0, err
	}

	updatedUser, err := r.GetById(ctx, id)
	if err != nil {
		return 0, err
	}
	return updatedUser.Rating, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT ` + userFields + ` FROM users WHERE email = $1`
	user := &dto.User{}
	err := r.db.Get(user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &entity.User{ID: user.ID, Username: user.Username, Email: user.Email, Rating: user.Rating, Role: user.Role, IsActive: user.IsActive, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT ` + userFields + ` FROM users WHERE username = $1`
	user := &dto.User{}
	err := r.db.Get(user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &entity.User{ID: user.ID, Username: user.Username, Email: user.Email, Rating: user.Rating, Role: user.Role, IsActive: user.IsActive, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}, nil
}

func (r *userRepository) CheckUserExists(ctx context.Context, id string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE id = $1`
	var count int
	err := r.db.Get(&count, query, id)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) GetPasswordById(ctx context.Context, id string) (string, error) {
	query := `SELECT password FROM users WHERE id = $1`
	var password string
	err := r.db.Get(&password, query, id)
	if err != nil {
		return "", err
	}
	return password, nil
}
