package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	application "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/profile"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
	profileMocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/mocks"
)

func TestProfileService_UpdateProfile_InvalidAge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := profileMocks.NewMockProfileRepository(ctrl)
	service := application.NewProfileService(mockProfileRepo, nil, nil, nil)
	ctx := context.Background()

	profile := &entity.Profile{UserID: "1", Age: -100}
	updated, err := service.UpdateProfile(ctx, profile)
	assert.Nil(t, updated)
	assert.Equal(t, application.ErrInvalidAge, err)
}

func TestProfileService_UpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepo := profileMocks.NewMockProfileRepository(ctrl)
	service := application.NewProfileService(mockProfileRepo, nil, nil, nil)
	ctx := context.Background()

	profile := &entity.Profile{UserID: "1", Age: 20}
	mockProfileRepo.EXPECT().UpdateProfile(ctx, profile).Return(profile, nil)

	updated, err := service.UpdateProfile(ctx, profile)
	assert.NoError(t, err)
	assert.Equal(t, profile, updated)
}
