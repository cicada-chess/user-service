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

func TestUserService_UpdateInfo_ErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	newUser := &entity.User{
		ID:       "1",
		Username: "new_user",
		Email:    "",
		Password: "password",
	}

	mockRepo.EXPECT().CheckUserExists(ctx, newUser.ID).Return(false, nil)

	_, err := userService.UpdateInfo(ctx, newUser)
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserService_UpdateInfo_ErrPasswordTooShort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	newUser := &entity.User{
		ID:       "1",
		Username: "new_user",
		Email:    "",
		Password: "short",
	}

	mockRepo.EXPECT().CheckUserExists(ctx, newUser.ID).Return(true, nil)

	_, err := userService.UpdateInfo(ctx, newUser)
	assert.Equal(t, entity.ErrPasswordTooShort, err)
}

func TestUserService_UpdateInfo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	newUser := &entity.User{
		ID:       "1",
		Username: "new_user",
		Email:    "",
		Password: "password",
	}

	newUserWithoutPass := &entity.User{
		ID:       "1",
		Username: "new_user",
		Email:    "",
	}

	mockRepo.EXPECT().CheckUserExists(ctx, newUser.ID).Return(true, nil)
	mockRepo.EXPECT().UpdateInfo(ctx, newUser).Return(newUserWithoutPass, nil)

	updatedUser, err := userService.UpdateInfo(ctx, newUser)
	assert.Nil(t, err)
	assert.Equal(t, newUserWithoutPass, updatedUser)
}
