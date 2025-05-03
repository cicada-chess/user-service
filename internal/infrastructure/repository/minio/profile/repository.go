package profile

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
)

type ProfileStorage struct {
	client *minio.Client
	bucket string
	host   string
}

func NewProfileStorage(client *minio.Client, bucket string, host string) interfaces.ProfileStorage {
	return &ProfileStorage{
		client: client,
		bucket: bucket,
		host:   host,
	}
}

func (s *ProfileStorage) SaveAvatar(ctx context.Context, objectName string, reader io.Reader, contentType string) (string, error) {
	userMetadata := map[string]string{"x-amz-acl": "public-read"}
	opts := minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: userMetadata,
	}

	if _, err := s.client.PutObject(ctx, s.bucket, objectName, reader, -1, opts); err != nil {
		return "", err
	}

	directURL := fmt.Sprintf("%s/%s/%s", s.host, s.bucket, objectName)
	return directURL, nil
}
