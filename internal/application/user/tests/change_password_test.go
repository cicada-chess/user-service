package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_ChangePassword_ErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(false, nil)

	err := userService.ChangePassword(ctx, "1", "old_password", "new_password")
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserService_ChangePassword_ErrInvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockPass, _ := entity.HashPassword("old_password")

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockRepo.EXPECT().GetPasswordById(ctx, "1").Return(mockPass, nil)

	err := userService.ChangePassword(ctx, "1", "wrong_old_password", "new_password")
	assert.Equal(t, user.ErrInvalidPassword, err)
}

func TestUserService_ChangePassword_ErrPasswordTooShort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)

	err := userService.ChangePassword(ctx, "1", "old_password", "short")
	assert.Equal(t, entity.ErrPasswordTooShort, err)
}

func TestUserService_ChangePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockPass, _ := entity.HashPassword("old_password")

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockRepo.EXPECT().GetPasswordById(ctx, "1").Return(mockPass, nil)
	mockRepo.EXPECT().ChangePassword(ctx, "1", gomock.Any()).Return(nil)

	err := userService.ChangePassword(ctx, "1", "old_password", "new_password")
	assert.Nil(t, err)
}
