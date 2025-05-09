package tests

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	application "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/profile"
)

func TestProfileService_UploadAvatar_InvalidFileType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := application.NewProfileService(nil, nil, nil, nil)
	ctx := context.Background()

	file := &multipart.FileHeader{Filename: "avatar.gif", Size: 100}
	url, err := service.UploadAvatar(ctx, "1", file)
	assert.Empty(t, url)
	assert.Equal(t, application.ErrInvalidFileType, err)
}

func TestProfileService_UploadAvatar_FileTooLarge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := application.NewProfileService(nil, nil, nil, nil)
	ctx := context.Background()

	file := &multipart.FileHeader{Filename: "avatar.jpg", Size: 6 * 1024 * 1024}
	url, err := service.UploadAvatar(ctx, "1", file)
	assert.Empty(t, url)
	assert.Equal(t, application.ErrFileSizeTooLarge, err)
}
