package filestorage

import (
	"bytes"
	"context"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
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
	// sniff max 512 bytes untuk deteksi MIME
	var sniff [512]byte
	n, _ := io.ReadFull(data, sniff[:]) // jika file <512B, error bisa diabaikan
	detectedCT := http.DetectContentType(sniff[:n])

	// reset reader ke awal
	data = io.MultiReader(bytes.NewReader(sniff[:n]), data)

	// jika key belum punya ekstensi → tambahkan dari content-type
	key = ensureExtIfMissing(key, detectedCT)

	// tentukan Content-Type final
	mimeType := detectedCT
	if mimeType == "" {
		if mt := mime.TypeByExtension(filepath.Ext(key)); mt != "" {
			mimeType = mt
		} else {
			mimeType = "application/octet-stream"
		}
	}

	_, err := s.Client.PutObject(ctx, s.BucketName, key, data, -1, minio.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return "", err
	}

	// kembalikan key final (sudah plus ekstensi)
	return key, nil
}

func ensureExtIfMissing(key, contentType string) string {
	// kalau sudah ada ekstensi, biarkan
	if filepath.Ext(key) != "" {
		return key
	}

	// ambil hanya media type-nya (buang parameter ;charset=)
	if mt, _, err := mime.ParseMediaType(contentType); err == nil && mt != "" {
		// coba dari registry MIME
		if exts, _ := mime.ExtensionsByType(mt); len(exts) > 0 {
			return key + strings.ToLower(exts[0])
		}
		// fallback manual untuk tipe umum
		switch strings.ToLower(mt) {
		case "image/jpeg", "image/jpg":
			return key + ".jpg"
		case "image/png":
			return key + ".png"
		case "image/webp":
			return key + ".webp"
		case "image/gif":
			return key + ".gif"
		}
	}

	// fallback terakhir
	return key + ".jpg"
}

// GetURL → generate presigned URL
func (s *S3Storage) GetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	params := make(url.Values)
	params.Set("response-content-disposition", "inline")

	if key == "" {
		return "", nil
	}

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

	if key == "" {
		return nil, nil
	}

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
