package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_GetRating_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(false, nil)

	rating, err := userService.GetRating(ctx, "1")
	assert.Equal(t, user.ErrUserNotFound, err)
	assert.Equal(t, 0, rating)
}

func TestUserService_GetRating_ErrInvalidUUIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	expectedError := &pq.Error{Severity: "ERROR", Code: "22P02"}
	mockRepo.EXPECT().CheckUserExists(ctx, "invalid").Return(false, expectedError)

	rating, err := userService.GetRating(ctx, "invalid")
	assert.Equal(t, user.ErrInvalidUUIDFormat, err)
	assert.Equal(t, 0, rating)
}

func TestUserService_GetRating_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockRepo.EXPECT().GetRating(ctx, "1").Return(1500, nil)

	rating, err := userService.GetRating(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, 1500, rating)
}
