package minio

import "os"

type StorageConfig struct {
	Endpoint   string
	Host       string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

func GetStorageConfig() StorageConfig {
	return StorageConfig{
		Endpoint:   os.Getenv("MINIO_ENDPOINT"),
		Host:       os.Getenv("MINIO_HOST"),
		AccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		BucketName: os.Getenv("MINIO_BUCKET_NAME"),
		UseSSL:     os.Getenv("MINIO_USE_SSL") == "true",
	}
}
