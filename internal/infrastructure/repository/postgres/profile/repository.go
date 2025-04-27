package user

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/repository/dto"
)

type profileRepository struct {
	db *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) interfaces.ProfileRepository {
	return &profileRepository{
		db: db,
	}
}

func (r *profileRepository) CreateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	query := `INSERT INTO profiles (user_id, description, age, location, avatar_path, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query, profile.UserID, profile.Description, profile.Age, profile.Location, profile.AvatarPath, profile.CreatedAt, profile.UpdatedAt)
	if err != nil {
		return nil, err
	}
	dbProfile, err := r.GetByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}
	return dbProfile, nil
}

func (r *profileRepository) GetByUserID(ctx context.Context, userID string) (*entity.Profile, error) {
	query := `SELECT user_id, description, age, location, avatar_path, created_at, updated_at FROM profiles WHERE user_id = $1`

	profile := &dto.Profile{}
	err := r.db.Get(profile, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &entity.Profile{
		UserID:      profile.UserID,
		Description: profile.Description,
		Age:         profile.Age,
		Location:    profile.Location,
		AvatarPath:  profile.AvatarPath,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}, nil
}

func (r *profileRepository) UpdateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	query := `
		UPDATE profiles
		SET description = $1, age = $2, location = $3, avatar_path = $4, updated_at = NOW()
		WHERE user_id = $5`

	_, err := r.db.Exec(query, profile.Description, profile.Age, profile.Location, profile.AvatarPath, profile.UserID)
	if err != nil {
		return nil, err
	}

	dbProfile, err := r.GetByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}

	return dbProfile, nil
}
