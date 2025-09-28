package filestorage

import (
	"context"
	"errors"
	"io"
	"prakarsa-app/config"
	"time"
)

type FileStorage interface {
	Put(ctx context.Context, path string, content io.Reader) (string, error)
	GetURL(ctx context.Context, key string, expiry time.Duration) (string, error)
	GetFile(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)
}

func NewFileStorage(cfg *config.Config) (FileStorage, error) {
	switch cfg.StorageType {
	case "local":
		return NewLocalStorage(cfg.LocalPath), nil
	case "s3":
		return NewS3Storage(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket, cfg.S3UseSSL)
	default:
		return nil, errors.New("unsupported storage type: " + cfg.StorageType)
	}
}
