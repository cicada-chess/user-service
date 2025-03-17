package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_UpdatePasswordById_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	err := userService.UpdatePasswordById(ctx, "1", "short")
	assert.Error(t, err)
}

func TestUserService_UpdatePasswordById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	validPassword := "ValidPassword1"

	mockRepo.EXPECT().ChangePassword(ctx, "1", gomock.Any()).Return(nil)

	err := userService.UpdatePasswordById(ctx, "1", validPassword)
	assert.NoError(t, err)
}
