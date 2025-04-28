package minio

import "os"

type DBConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

func GetDBConfig() DBConfig {
	return DBConfig{
		Endpoint:   os.Getenv("MINIO_ENDPOINT"),
		AccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		BucketName: os.Getenv("MINIO_BUCKET_NAME"),
		UseSSL:     os.Getenv("MINIO_USE_SSL") == "true",
	}
}
