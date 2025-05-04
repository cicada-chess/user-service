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

func TestUserService_ToggleActive_ErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(false, nil)

	_, err := userService.ToggleActive(ctx, "1")
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserService_ToggleActive_ErrInvalidUUIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	expectedError := &pq.Error{Severity: "ERROR", Code: "22P02"}
	mockRepo.EXPECT().CheckUserExists(ctx, "invalid").Return(false, expectedError)

	_, err := userService.ToggleActive(ctx, "invalid")
	assert.Equal(t, user.ErrInvalidUUIDFormat, err)

}

func TestUserService_ToggleActive_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockRepo.EXPECT().ToggleActive(ctx, "1").Return(true, nil)

	isActive, err := userService.ToggleActive(ctx, "1")
	assert.Nil(t, err)
	assert.True(t, isActive)
}
