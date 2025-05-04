package config

import "os"

type Config struct {
	Storage StorageConfig
	DB      DBConfig
}

type StorageConfig struct {
	Endpoint   string
	Host       string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func ReadConfig() (*Config, error) {
	return &Config{
		Storage: StorageConfig{
			Endpoint:   os.Getenv("MINIO_ENDPOINT"),
			Host:       os.Getenv("MINIO_HOST"),
			AccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
			SecretKey:  os.Getenv("MINIO_SECRET_KEY"),
			BucketName: os.Getenv("MINIO_BUCKET_NAME"),
			UseSSL:     os.Getenv("MINIO_USE_SSL") == "true",
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("SSL_MODE"),
		},
	}, nil
}
