package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	mocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
)

func TestUserService_Create_ErrUsernameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	existingUser := &entity.User{
		Username: "existing_user",
		Email:    "existing@example.com",
	}
	mockRepo.EXPECT().GetByEmail(ctx, "new@example.com").Return(nil, sql.ErrNoRows)
	mockRepo.EXPECT().GetByUsername(ctx, "existing_user").Return(existingUser, nil)

	newUser := &entity.User{
		Username: "existing_user",
		Email:    "new@example.com",
	}

	user, err := userService.Create(context.Background(), newUser)
	assert.Nil(t, user)
	assert.NotNil(t, err)
}

func TestUserService_Create_ErrEmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	existingUser := &entity.User{
		Username: "existing_user",
		Email:    "existing@example.com",
	}

	mockRepo.EXPECT().GetByEmail(ctx, "existing@example.com").Return(existingUser, nil)

	newUser := &entity.User{
		Username: "new_user",
		Email:    "existing@example.com",
	}

	user, err := userService.Create(context.Background(), newUser)
	assert.Nil(t, user)
	assert.NotNil(t, err)
}

func TestUserService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)
	ctx := context.Background()

	newUser := &entity.User{
		Username: "new_user",
		Email:    "new@example.com",
	}

	mockRepo.EXPECT().GetByEmail(ctx, "new@example.com").Return(nil, sql.ErrNoRows)
	mockRepo.EXPECT().GetByUsername(ctx, "new_user").Return(nil, sql.ErrNoRows)
	mockRepo.EXPECT().Create(ctx, newUser).Return(newUser, nil)

	createdUser, err := userService.Create(ctx, newUser)
	assert.NotNil(t, createdUser)
	assert.Nil(t, err)
	assert.Equal(t, "new_user", createdUser.Username)
	assert.Equal(t, "new@example.com", createdUser.Email)
}
