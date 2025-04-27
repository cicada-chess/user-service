package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	pb "gitlab.mai.ru/cicada-chess/backend/auth-service/pkg/auth"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
	userInterfaces "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidFileType       = errors.New("invalid file type")
	ErrFileSizeTooLarge      = errors.New("file size too large")
	ErrInvalidAge            = errors.New("invalid age")
	ErrTokenInvalidOrExpired = errors.New("token is invalid or expired")
	ErrInternalServer        = errors.New("internal server error")
)

type profileService struct {
	profileRepo interfaces.ProfileRepository
	userRepo    userInterfaces.UserRepository
	client      pb.AuthServiceClient
}

// NewProfileService создает новый сервис для работы с профилями
func NewProfileService(profileRepo interfaces.ProfileRepository, userRepo userInterfaces.UserRepository, client pb.AuthServiceClient) interfaces.ProfileService {
	return &profileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
		client:      client,
	}
}

// GetProfile возвращает профиль пользователя
func (s *profileService) GetProfile(ctx context.Context, userID string) (*entity.Profile, error) {
	// Проверяем, существует ли пользователь
	exists, err := s.userRepo.CheckUserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	// Получаем профиль
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Если профиля нет, возвращаем пустой профиль
	if profile == nil {
		profile = &entity.Profile{
			UserID: userID,
		}
	}

	return profile, nil
}

// CreateOrUpdateProfile создает или обновляет профиль пользователя
func (s *profileService) CreateOrUpdateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	// Проверяем, существует ли пользователь
	exists, err := s.userRepo.CheckUserExists(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	// Валидируем профиль
	if !profile.IsValidAge() {
		return nil, ErrInvalidAge
	}

	// Проверяем, существует ли профиль
	existingProfile, err := s.profileRepo.GetByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}

	if existingProfile == nil {
		// Если профиля нет, создаем новый
		profile.CreatedAt = time.Now()
		profile.UpdatedAt = time.Now()
		return s.profileRepo.Create(ctx, profile)
	}

	// Если профиль существует, обновляем его
	existingProfile.Description = profile.Description
	existingProfile.Age = profile.Age
	existingProfile.Location = profile.Location
	existingProfile.UpdatedAt = time.Now()

	// Сохраняем путь к аватару, если был новый
	if profile.AvatarPath != "" {
		existingProfile.AvatarPath = profile.AvatarPath
	}

	return s.profileRepo.Update(ctx, existingProfile)
}

// UploadAvatar загружает аватар пользователя
func (s *profileService) UploadAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	// Проверяем, существует ли пользователь
	exists, err := s.userRepo.CheckUserExists(ctx, userID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrUserNotFound
	}

	// Проверка типа файла
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return "", ErrInvalidFileType
	}

	// Проверка размера файла (максимум 5MB)
	if file.Size > 5*1024*1024 {
		return "", ErrFileSizeTooLarge
	}

	// Открываем загруженный файл
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Считываем данные файла
	fileData, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Сохраняем аватар в хранилище через репозиторий
	avatarPath, err := s.profileRepo.SaveAvatar(ctx, userID, fileData, ext)
	if err != nil {
		return "", fmt.Errorf("failed to save avatar: %w", err)
	}

	// Обновляем профиль с новым путем к аватару
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	if profile != nil {
		profile.AvatarPath = avatarPath
		profile.UpdatedAt = time.Now()
		_, err = s.profileRepo.Update(ctx, profile)
		if err != nil {
			return "", err
		}
	} else {
		// Если профиля нет, создаем новый
		newProfile := &entity.Profile{
			UserID:     userID,
			AvatarPath: avatarPath,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_, err = s.profileRepo.Create(ctx, newProfile)
		if err != nil {
			return "", err
		}
	}

	return avatarPath, nil
}

func (s *profileService) GetUserIDFromToken(ctx context.Context, tokenHeader string) (string, error) {
	accessToken := strings.TrimPrefix(tokenHeader, "Bearer ")
	req := &pb.ValidateTokenRequest{Token: tokenHeader}
	_, err := s.client.ValidateToken(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.PermissionDenied:
				return "", ErrTokenInvalidOrExpired
			default:
				return "", err
			}
		}
	}

	token, _ := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalidOrExpired
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)
	userId, _ := claims["user_id"].(string)

	return userId, nil
}
