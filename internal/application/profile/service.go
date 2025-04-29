package user

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
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
	ErrInvalidUUIDFormat     = errors.New("invalid uuid format")
)

type profileService struct {
	profileRepo    interfaces.ProfileRepository
	userRepo       userInterfaces.UserRepository
	profileStorage interfaces.ProfileStorage
	client         pb.AuthServiceClient
}

func NewProfileService(profileRepo interfaces.ProfileRepository, userRepo userInterfaces.UserRepository, profileStorage interfaces.ProfileStorage, client pb.AuthServiceClient) interfaces.ProfileService {
	return &profileService{
		profileRepo:    profileRepo,
		userRepo:       userRepo,
		profileStorage: profileStorage,
		client:         client,
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
		UserID:      userID,
		Age:         -1,
		Description: "",
		Location:    "",
		AvatarURL:   "",
	}

	createdProfile, err := s.profileRepo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return createdProfile, nil
}

func (s *profileService) GetProfile(ctx context.Context, userID string) (*entity.Profile, error) {
	exists, err := s.userRepo.CheckUserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	exists, err = s.profileRepo.CheckProfileExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrProfileNotFound
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "22P02" {
			return nil, ErrInvalidUUIDFormat
		}
		return nil, err
	}

	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	if !profile.IsValidAge() && profile.Age != -1 {
		return nil, ErrInvalidAge
	}

	updatedProfile, err := s.profileRepo.UpdateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return updatedProfile, nil
}

func (s *profileService) UploadAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", ErrInvalidFileType
	}
	if file.Size > 5*1024*1024 {
		return "", ErrFileSizeTooLarge
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	objectName := fmt.Sprintf("%s%s", userID, ext)
	contentType := file.Header.Get("Content-Type")

	url, err := s.profileStorage.SaveAvatar(ctx, objectName, src, contentType)
	if err != nil {
		return "", err
	}
	return url, nil
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
