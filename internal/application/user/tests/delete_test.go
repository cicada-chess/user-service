package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_Delete_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(false, nil)

	err := userService.Delete(ctx, "1")
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockRepo.EXPECT().Delete(ctx, "1").Return(nil)

	err := userService.Delete(ctx, "1")

	assert.Nil(t, err)
}
