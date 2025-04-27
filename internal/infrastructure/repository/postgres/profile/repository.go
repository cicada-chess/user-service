package user

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
)

type profileRepository struct {
	db *sqlx.DB
}

// NewProfileRepository создает новый репозиторий профилей
func NewProfileRepository(db *sqlx.DB) interfaces.ProfileRepository {
	return &profileRepository{
		db: db,
	}
}

// Create создает новый профиль
func (r *profileRepository) Create(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	query := `
		INSERT INTO profiles (user_id, description, age, location, avatar_path, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id, description, age, location, avatar_path, created_at, updated_at
	`

	row := r.db.QueryRowContext(
		ctx, query, profile.UserID, profile.Description, profile.Age, profile.Location,
		profile.AvatarPath, profile.CreatedAt, profile.UpdatedAt,
	)

	var result entity.Profile
	err := row.Scan(&result.UserID, &result.Description, &result.Age,
		&result.Location, &result.AvatarPath, &result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// Обработка дубликатов (например, если профиль с таким user_id уже существует)
			return nil, fmt.Errorf("profile for this user already exists")
		}
		return nil, err
	}

	return &result, nil
}

// GetByUserID возвращает профиль по ID пользователя
func (r *profileRepository) GetByUserID(ctx context.Context, userID string) (*entity.Profile, error) {
	query := `
		SELECT user_id, description, age, location, avatar_path, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, userID)

	var profile entity.Profile
	err := row.Scan(&profile.UserID, &profile.Description, &profile.Age,
		&profile.Location, &profile.AvatarPath, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Профиль не найден
		}
		return nil, err
	}

	return &profile, nil
}

// Update обновляет существующий профиль
func (r *profileRepository) Update(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	query := `
		UPDATE profiles
		SET description = $1, age = $2, location = $3, avatar_path = $4, updated_at = $5
		WHERE user_id = $6
		RETURNING user_id, description, age, location, avatar_path, created_at, updated_at
	`

	row := r.db.QueryRowContext(
		ctx, query,
		profile.Description, profile.Age, profile.Location, profile.AvatarPath,
		time.Now(), profile.UserID,
	)

	var result entity.Profile
	err := row.Scan(&result.UserID, &result.Description, &result.Age,
		&result.Location, &result.AvatarPath, &result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

// SaveAvatar сохраняет аватар на диск и возвращает путь
func (r *profileRepository) SaveAvatar(ctx context.Context, userID string, avatarData []byte, fileExt string) (string, error) {
	avatarsDir := "uploads/avatars"
	if err := os.MkdirAll(avatarsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Генерируем уникальное имя файла
	avatarFileName := fmt.Sprintf("%s-%s%s", userID, uuid.New().String()[:8], fileExt)
	avatarPath := filepath.Join(avatarsDir, avatarFileName)

	if err := os.WriteFile(avatarPath, avatarData, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return avatarPath, nil
}
