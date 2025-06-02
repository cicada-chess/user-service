package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	notificationMocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/notification/mocks"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_Create_ErrUsernameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	existingUser := &entity.User{
		Username: "existing_user",
		Email:    "existing@example.com",
	}
	mockRepo.EXPECT().GetByEmail(ctx, "new@example.com").Return(nil, nil)
	mockRepo.EXPECT().GetByUsername(ctx, "existing_user").Return(existingUser, nil)

	newUser := &entity.User{
		Username: "existing_user",
		Email:    "new@example.com",
		Password: "password",
	}

	createdUser, err := userService.Create(context.Background(), newUser)
	assert.Nil(t, createdUser)
	assert.Equal(t, user.ErrUsernameExists, err)
}

func TestUserService_Create_ErrEmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	existingUser := &entity.User{
		Username: "existing_user",
		Email:    "existing@example.com",
	}

	mockRepo.EXPECT().GetByEmail(ctx, "existing@example.com").Return(existingUser, nil)

	newUser := &entity.User{
		Username: "new_user",
		Email:    "existing@example.com",
		Password: "password",
	}

	createdUser, err := userService.Create(context.Background(), newUser)
	assert.Nil(t, createdUser)
	assert.Equal(t, user.ErrEmailExists, err)
}

func TestUserService_Create_ErrPasswordTooShort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	newUser := &entity.User{
		Username: "new_user",
		Email:    "new@example.com",
		Password: "pass",
	}

	mockRepo.EXPECT().GetByEmail(ctx, "new@example.com").Return(nil, nil)
	mockRepo.EXPECT().GetByUsername(ctx, "new_user").Return(nil, nil)

	createdUser, err := userService.Create(context.Background(), newUser)
	assert.Nil(t, createdUser)
	assert.Equal(t, entity.ErrPasswordTooShort, err)
}

func TestUserService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockNotificationSender := notificationMocks.NewMockNotificationSender(ctrl)
	userService := user.NewUserService(mockRepo, mockNotificationSender)

	newUser := &entity.User{
		Username: "new_user",
		Email:    "new@example.com",
		Password: "password",
	}

	mockRepo.EXPECT().GetByEmail(ctx, "new@example.com").Return(nil, nil)
	mockRepo.EXPECT().GetByUsername(ctx, "new_user").Return(nil, nil)
	mockRepo.EXPECT().Create(ctx, newUser).Return(newUser, nil)
	mockNotificationSender.EXPECT().SendAccountConfirmation(ctx, newUser.Email, newUser.Username, gomock.Any()).Return(nil)

	createdUser, err := userService.Create(ctx, newUser)
	assert.NotNil(t, createdUser)
	assert.Nil(t, err)
	assert.Equal(t, newUser.Username, createdUser.Username)
	assert.Equal(t, newUser.Email, createdUser.Email)
}
