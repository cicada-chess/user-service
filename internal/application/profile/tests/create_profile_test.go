package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	application "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/profile"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
	profileMocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/mocks"
	userMocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestProfileService_CreateProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := profileMocks.NewMockProfileRepository(ctrl)
	mockUserRepo := userMocks.NewMockUserRepository(ctrl)
	service := application.NewProfileService(mockProfileRepo, mockUserRepo, nil, nil)
	ctx := context.Background()

	mockUserRepo.EXPECT().CheckUserExists(ctx, "1").Return(false, nil)

	profile, err := service.CreateProfile(ctx, "1")
	assert.Nil(t, profile)
	assert.Equal(t, application.ErrUserNotFound, err)
}

func TestProfileService_CreateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := profileMocks.NewMockProfileRepository(ctrl)
	mockUserRepo := userMocks.NewMockUserRepository(ctrl)
	service := application.NewProfileService(mockProfileRepo, mockUserRepo, nil, nil)
	ctx := context.Background()

	mockUserRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockUserRepo.EXPECT().GetUsernameByUserID(ctx, "1").Return("testuser", nil)
	expectedProfile := &entity.Profile{UserID: "1", Username: "testuser", Age: -1}
	mockProfileRepo.EXPECT().CreateProfile(ctx, gomock.Any()).Return(expectedProfile, nil)

	profile, err := service.CreateProfile(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, expectedProfile, profile)
}
