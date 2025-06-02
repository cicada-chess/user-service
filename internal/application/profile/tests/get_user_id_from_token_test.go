package tests

import (
	"context"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	pb "gitlab.mai.ru/cicada-chess/backend/auth-service/pkg/auth"
	application "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/profile"
	authMocks "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/auth/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProfileService_GetUserIDFromToken_PermissionDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := authMocks.NewMockAuthServiceClient(ctrl)
	service := application.NewProfileService(nil, nil, nil, mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ValidateToken(ctx, &pb.ValidateTokenRequest{Token: "Bearer badtoken"}).
		Return(nil, status.Error(codes.PermissionDenied, "invalid token"))

	userID, err := service.GetUserIDFromToken(ctx, "Bearer badtoken")
	assert.Empty(t, userID)
	assert.Equal(t, application.ErrTokenInvalidOrExpired, err)
}

func TestProfileService_GetUserIDFromToken_OtherGRPCError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := authMocks.NewMockAuthServiceClient(ctrl)
	service := application.NewProfileService(nil, nil, nil, mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ValidateToken(ctx, &pb.ValidateTokenRequest{Token: "Bearer sometoken"}).
		Return(nil, status.Error(codes.Internal, "internal error"))

	userID, err := service.GetUserIDFromToken(ctx, "Bearer sometoken")
	assert.Empty(t, userID)
	assert.Error(t, err)
}

func TestProfileService_GetUserIDFromToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := authMocks.NewMockAuthServiceClient(ctrl)
	service := application.NewProfileService(nil, nil, nil, mockClient)
	ctx := context.Background()

	// Валидный токен с user_id
	secret := "testsecret"
	_ = os.Setenv("SECRET_KEY", secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "12345",
	})
	tokenString, _ := token.SignedString([]byte(secret))

	mockClient.EXPECT().ValidateToken(ctx, &pb.ValidateTokenRequest{Token: "Bearer " + tokenString}).Return(nil, nil)

	userID, err := service.GetUserIDFromToken(ctx, "Bearer "+tokenString)
	assert.NoError(t, err)
	assert.Equal(t, "12345", userID)
}
