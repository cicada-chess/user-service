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

func TestUserService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := user.NewUserService(mockRepo)

	ctx := context.Background()
	page := "1"
	limit := "10"
	search := "test"
	sortBy := "username"
	order := "asc"

	expectedUsers := []*entity.User{
		{ID: "1", Username: "testuser1", Email: "test1@example.com"},
		{ID: "2", Username: "testuser2", Email: "test2@example.com"},
	}

	mockRepo.EXPECT().GetAll(ctx, page, limit, search, sortBy, order).Return(expectedUsers, nil)

	users, err := userService.GetAll(ctx, page, limit, search, sortBy, order)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}
