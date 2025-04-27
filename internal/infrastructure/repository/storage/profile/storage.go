package profile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	interfaces "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
)

type profileStorage struct {
	Directory string
}

func NewProfileStorage(directory string) interfaces.ProfileStorage {
	return &profileStorage{Directory: directory}
}

func (s *profileStorage) SaveAvatar(ctx context.Context, userID string, avatarData []byte, fileExt string) (string, error) {
	avatarsDir := s.Directory
	if err := os.MkdirAll(avatarsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	avatarFileName := fmt.Sprintf("%s-%s%s", userID, uuid.New().String()[:8], fileExt)
	avatarPath := filepath.Join(avatarsDir, avatarFileName)

	if err := os.WriteFile(avatarPath, avatarData, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return avatarPath, nil
}
