package profile

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
)

type ProfileStorage struct {
	client *minio.Client
	bucket string
}

func NewProfileStorage(client *minio.Client, bucket string) interfaces.ProfileStorage {
	return &ProfileStorage{
		client: client,
		bucket: bucket,
	}
}

func (s *ProfileStorage) SaveAvatar(ctx context.Context, objectName string, reader io.Reader, contentType string) (string, error) {
	opts := minio.PutObjectOptions{ContentType: contentType}
	if _, err := s.client.PutObject(ctx, s.bucket, objectName, reader, -1, opts); err != nil {
		return "", err
	}
	presigned, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, time.Hour*24*7, url.Values{})
	if err != nil {
		return "", err
	}
	return presigned.String(), nil
}
