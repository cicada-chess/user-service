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
	ErrProfileNotFound       = errors.New("profile not found")
)

type profileService struct {
	profileRepo interfaces.ProfileRepository
	userRepo    userInterfaces.UserRepository
	client      pb.AuthServiceClient
}

func NewProfileService(profileRepo interfaces.ProfileRepository, userRepo userInterfaces.UserRepository, client pb.AuthServiceClient) interfaces.ProfileService {
	return &profileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
		client:      client,
	}
}

func (s *profileService) CreateProfile(ctx context.Context, userID string) (*entity.Profile, error) {
	exists, err := s.userRepo.CheckUserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	profile := &entity.Profile{
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.profileRepo.CreateProfile(ctx, profile)
}

func (s *profileService) GetProfile(ctx context.Context, userID string) (*entity.Profile, error) {
	exists, err := s.userRepo.CheckUserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		profile = &entity.Profile{
			UserID: userID,
		}
	}

	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	exists, err := s.userRepo.CheckUserExists(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	if !profile.IsValidAge() {
		return nil, ErrInvalidAge
	}

	existingProfile, err := s.profileRepo.GetByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}

	if existingProfile == nil {
		return s.profileRepo.CreateProfile(ctx, profile)
	}

	return s.profileRepo.UpdateProfile(ctx, profile)
}

func (s *profileService) UploadAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	exists, err := s.userRepo.CheckUserExists(ctx, userID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrUserNotFound
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", ErrInvalidFileType
	}

	if file.Size > 5*1024*1024 {
		return "", ErrFileSizeTooLarge
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	if profile != nil && profile.AvatarPath != "" {
		if err := os.Remove(profile.AvatarPath); err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to remove old avatar: %w", err)
		}
	}

	avatarPath, err := s.profileRepo.SaveAvatar(ctx, userID, fileData, ext)
	if err != nil {
		return "", fmt.Errorf("failed to save avatar: %w", err)
	}

	if profile != nil {
		profile.AvatarPath = avatarPath
		_, err = s.profileRepo.UpdateProfile(ctx, profile)
		if err != nil {
			return "", err
		}
	} else {
		newProfile := &entity.Profile{
			UserID:     userID,
			AvatarPath: avatarPath,
		}
		_, err = s.profileRepo.CreateProfile(ctx, newProfile)
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
