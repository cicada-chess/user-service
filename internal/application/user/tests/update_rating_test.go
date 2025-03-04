package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_UpdateRating_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(false, nil)

	rating, err := userService.UpdateRating(ctx, "1", 100)
	assert.Equal(t, user.ErrUserNotFound, err)
	assert.Equal(t, 0, rating)
}

func TestUserService_UpdateRating_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockRepo.EXPECT().UpdateRating(ctx, "1", 100).Return(1600, nil)

	rating, err := userService.UpdateRating(ctx, "1", 100)
	assert.NoError(t, err)
	assert.Equal(t, 1600, rating)
}
