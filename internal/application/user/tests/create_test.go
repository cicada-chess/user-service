package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	service "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	mock_interfaces "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/mocks"
	"go.uber.org/mock/gomock"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_interfaces.NewMockUserRepository(ctrl)

	userID := uuid.UUID{}.String()
	service := service.NewUserService(repoMock)

	newUser := &entity.User{ID: userID, Username: "Anonymous", Email: "K1M3w@example.com", Password: "password", Role: 1}

	ctx := context.Background()
	repoMock.EXPECT().
		Create(ctx, newUser).
		Return(userID, nil).
		Times(1)

	_, err := service.Create(ctx, newUser)

	assert.NoError(t, err)
}

func TestUserService_Create_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_interfaces.NewMockUserRepository(ctrl)

	service := service.NewUserService(repoMock)

	newUser := &entity.User{Username: "Anonymous", Email: "K1M3w@example.com", Password: "password", Role: 1}
	expectedError := errors.New("repository error")

	ctx := context.Background()
	repoMock.EXPECT().
		Create(ctx, newUser).
		Return("", expectedError).
		Times(1)

	id, err := service.Create(ctx, newUser)

	assert.Error(t, err)
	assert.Equal(t, "", id)
	assert.Equal(t, expectedError, err)
}
