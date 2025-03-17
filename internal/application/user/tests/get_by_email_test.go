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

func TestUserService_GetByEmail_ErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().GetByEmail(ctx, "example@example.com").Return(nil, nil)

	_, err := userService.GetUserByEmail(ctx, "example@example.com")
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserService_GetByEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	expectedUser := &entity.User{
		ID:       "1",
		Username: "testuser",
		Email:    "example@example.com",
		Password: "password",
	}

	mockRepo.EXPECT().GetByEmail(ctx, "example@example.com").Return(expectedUser, nil)
	mockRepo.EXPECT().GetPasswordById(ctx, "1").Return("password", nil)
	user, err := userService.GetUserByEmail(ctx, "example@example.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}
