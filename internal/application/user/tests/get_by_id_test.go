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

func TestUserService_GetById_ErrUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetById(ctx, "1").Return(nil, nil)

	userService := user.NewUserService(mockRepo, nil)

	_, err := userService.GetById(ctx, "1")
	assert.Equal(t, user.ErrUserNotFound, err)
}

func TestUserService_GetById_ErrInvalidUUIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo, nil)

	expectedError := &pq.Error{Severity: "ERROR", Code: "22P02"}
	mockRepo.EXPECT().GetById(ctx, "invalid").Return(nil, expectedError)

	_, err := userService.GetById(ctx, "invalid")
	assert.Equal(t, user.ErrInvalidUUIDFormat, err)
}

func TestUserService_GetById_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	ctx := context.Background()
	userService := user.NewUserService(mockRepo, nil)

	mockRepo.EXPECT().GetById(ctx, "1").Return(&entity.User{ID: "1"}, nil)

	user, err := userService.GetById(ctx, "1")
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)

}
