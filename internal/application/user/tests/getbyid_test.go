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

func TestUserService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_interfaces.NewMockUserRepository(ctrl)

	service := service.NewUserService(repoMock)

	userID := uuid.UUID{}.String()
	expectedUser := &entity.User{ID: userID, Username: "Anonymous", Email: "K1M3w@example.com", Password: "password", Role: 1}

	ctx := context.Background()
	repoMock.EXPECT().
		GetById(ctx, userID).
		Return(expectedUser, nil).
		Times(1)

	user, err := service.GetById(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_GetById_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock_interfaces.NewMockUserRepository(ctrl)

	service := service.NewUserService(repoMock)

	userID := uuid.UUID{}.String()

	ctx := context.Background()
	repoMock.EXPECT().
		GetById(ctx, userID).
		Return(nil, errors.New("error")).
		Times(1)

	user, err := service.GetById(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)
}
