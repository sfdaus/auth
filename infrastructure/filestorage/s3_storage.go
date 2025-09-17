package filestorage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	Client     *minio.Client
	BucketName string
}

func NewS3Storage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*S3Storage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		Client:     minioClient,
		BucketName: bucketName,
	}, nil
}

func (s *S3Storage) Put(ctx context.Context, key string, data io.Reader) (string, error) {
	ext := filepath.Ext(key) // ambil extension dari nama file
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// fallback → coba sniff dari isi file
		buf := make([]byte, 512)
		n, _ := data.Read(buf)
		mimeType = http.DetectContentType(buf[:n])

		// reset reader ke awal karena udah di-read
		data = io.MultiReader(bytes.NewReader(buf[:n]), data)
	}

	_, err := s.Client.PutObject(ctx, s.BucketName, key, data, -1, minio.PutObjectOptions{ContentType: mimeType})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s", key)
	return url, nil
}

// GetURL → generate presigned URL
func (s *S3Storage) GetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	params := make(url.Values)
	params.Set("response-content-disposition", "inline")

	req, err := s.Client.PresignedGetObject(ctx, s.BucketName, key, expiry, params)
	if err != nil {
		return "", err
	}
	return req.String(), nil
}

// GetFile → download object content
func (s *S3Storage) GetFile(ctx context.Context, key string) ([]byte, error) {
	obj, err := s.Client.GetObject(ctx, s.BucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	return s.Client.RemoveObject(ctx, s.BucketName, key, minio.RemoveObjectOptions{})
}

func (s *S3Storage) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.Client.StatObject(ctx, s.BucketName, path, minio.StatObjectOptions{})
	if err != nil {
		// kalau object tidak ada → StatObject error
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
