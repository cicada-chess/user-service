package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"

	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_UpdatePasswordById_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	err := userService.UpdatePasswordById(ctx, "1", "short")
	assert.Equal(t, entity.ErrPasswordTooShort, err)
}

func TestUserService_UpdatePasswordById_ErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	mockRepo.EXPECT().CheckUserExists(ctx, "2").Return(false, nil)
	err := userService.UpdatePasswordById(ctx, "2", "ValidPassword1")
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserService_UpdatePasswordById_ErrInvalidUUIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	expectedError := &pq.Error{Severity: "ERROR", Code: "22P02"}
	mockRepo.EXPECT().CheckUserExists(ctx, "invalid").Return(false, expectedError)
	err := userService.UpdatePasswordById(ctx, "invalid", "ValidPassword1")
	assert.Equal(t, user.ErrInvalidUUIDFormat, err)
}

func TestUserService_UpdatePasswordById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	validPassword := "ValidPassword1"
	mockRepo.EXPECT().CheckUserExists(ctx, "1").Return(true, nil)
	mockRepo.EXPECT().ChangePassword(ctx, "1", gomock.Any()).Return(nil)

	err := userService.UpdatePasswordById(ctx, "1", validPassword)
	assert.NoError(t, err)
}
