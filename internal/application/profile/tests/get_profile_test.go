package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	application "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/profile"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
	profileMocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/mocks"
	userMocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestProfileService_GetProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := profileMocks.NewMockProfileRepository(ctrl)
	mockUserRepo := userMocks.NewMockUserRepository(ctrl)
	service := application.NewProfileService(mockProfileRepo, mockUserRepo, nil, nil)
	ctx := context.Background()

	mockUserRepo.EXPECT().CheckUserExists(ctx, "1").Return(false, nil)

	profile, err := service.GetProfile(ctx, "1")
	assert.Nil(t, profile)
	assert.Equal(t, application.ErrUserNotFound, err)
}

func TestProfileService_GetProfile_ProfileNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := profileMocks.NewMockProfileRepository(ctrl)
	mockUserRepo := userMocks.NewMockUserRepository(ctrl)
	service := application.NewProfileService(mockProfileRepo, mockUserRepo, nil, nil)
	ctx := context.Background()

	mockUserRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockProfileRepo.EXPECT().CheckProfileExists(ctx, "1").Return(false, nil)

	profile, err := service.GetProfile(ctx, "1")
	assert.Nil(t, profile)
	assert.Equal(t, application.ErrProfileNotFound, err)
}

func TestProfileService_GetProfile_ErrInvalidUUIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := profileMocks.NewMockProfileRepository(ctrl)
	mockUserRepo := userMocks.NewMockUserRepository(ctrl)
	service := application.NewProfileService(mockProfileRepo, mockUserRepo, nil, nil)
	ctx := context.Background()

	mockUserRepo.EXPECT().CheckUserExists(ctx, "bad").Return(true, nil)
	mockProfileRepo.EXPECT().CheckProfileExists(ctx, "bad").Return(true, nil)
	mockProfileRepo.EXPECT().GetByUserID(ctx, "bad").Return(nil, &pq.Error{Severity: "ERROR", Code: "22P02"})

	profile, err := service.GetProfile(ctx, "bad")
	assert.Nil(t, profile)
	assert.Equal(t, application.ErrInvalidUUIDFormat, err)
}

func TestProfileService_GetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := profileMocks.NewMockProfileRepository(ctrl)
	mockUserRepo := userMocks.NewMockUserRepository(ctrl)
	service := application.NewProfileService(mockProfileRepo, mockUserRepo, nil, nil)
	ctx := context.Background()

	expectedProfile := &entity.Profile{UserID: "1", Username: "testuser"}
	mockUserRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockProfileRepo.EXPECT().CheckProfileExists(ctx, "1").Return(true, nil)
	mockProfileRepo.EXPECT().GetByUserID(ctx, "1").Return(expectedProfile, nil)

	profile, err := service.GetProfile(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile, profile)
}
