package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	notificationMocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/notification/mocks"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_ForgotPassword_ErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	mockRepo.EXPECT().GetByEmail(ctx, "notfound@example.com").Return(nil, nil)

	err := userService.ForgotPassword(ctx, "notfound@example.com")
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserService_ForgotPassword_ErrInvalidUUIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)
	ctx := context.Background()

	expectedError := &pq.Error{Severity: "ERROR", Code: "22P02"}
	mockRepo.EXPECT().GetByEmail(ctx, "bad@example.com").Return(nil, expectedError)

	err := userService.ForgotPassword(ctx, "bad@example.com")
	assert.Equal(t, user.ErrInvalidUUIDFormat, err)
}

func TestUserService_ForgotPassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockNotificationSender := notificationMocks.NewMockNotificationSender(ctrl)
	userService := user.NewUserService(mockRepo, mockNotificationSender)
	ctx := context.Background()

	testUser := &entity.User{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
	}

	mockRepo.EXPECT().GetByEmail(ctx, "test@example.com").Return(testUser, nil)
	mockNotificationSender.EXPECT().SendPasswordReset(ctx, testUser.Email, testUser.Username, gomock.Any()).Return(nil)

	err := userService.ForgotPassword(ctx, "test@example.com")
	assert.Nil(t, err)
}
