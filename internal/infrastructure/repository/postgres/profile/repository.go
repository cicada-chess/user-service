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
	query := `INSERT INTO profiles (user_id, username, age, description, location, avatar_url) VALUES ($1, $2, -1, '', '', '')`

	_, err := r.db.Exec(query, profile.UserID, profile.Username)
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
	query := `SELECT * FROM profiles WHERE user_id = $1`

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
		Username:    profile.Username,
		Description: profile.Description,
		Age:         profile.Age,
		Location:    profile.Location,
		AvatarURL:   profile.AvatarURL,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}, nil
}

func (r *profileRepository) UpdateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	query := `UPDATE profiles SET description = $2, age = $3, location = $4, avatar_url = $5, updated_at = NOW() WHERE user_id = $1`

	_, err := r.db.Exec(query, profile.UserID, profile.Description, profile.Age, profile.Location, profile.AvatarURL)
	if err != nil {
		return nil, err
	}

	dbProfile, err := r.GetByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}

	return dbProfile, nil
}

func (r *profileRepository) CheckProfileExists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT COUNT(*) FROM profiles WHERE user_id = $1`
	var count int
	err := r.db.Get(&count, query, userID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
