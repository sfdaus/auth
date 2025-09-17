package filestorage

import (
	"bytes"
	"context"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"prakarsa-app/config"
	"time"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

func (l *LocalStorage) Put(ctx context.Context, key string, data io.Reader) (string, error) {
	ext := filepath.Ext(key)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		buf := make([]byte, 512)
		n, _ := data.Read(buf)
		mimeType = http.DetectContentType(buf[:n])
		data = io.MultiReader(bytes.NewReader(buf[:n]), data)
	}

	// pastikan folder ada
	fullPath := filepath.Join(l.BasePath, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return "", err
	}

	// return URL publik (bisa di-serve lewat HTTP server)
	return key, nil
}

// func (l *LocalStorage) Get(path string) (io.ReadCloser, error) {
// 	fullPath := filepath.Join(l.BasePath, path)
// 	return os.Open(fullPath)
// }

func (l *LocalStorage) GetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	return config.LoadConfig().PublicURL + l.BasePath + key, nil
}

func (l *LocalStorage) GetFile(ctx context.Context, key string) ([]byte, error) {
	fullPath := filepath.Join(l.BasePath, key)
	return os.ReadFile(fullPath)
}

func (l *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(l.BasePath, path)
	return os.Remove(fullPath)
}

func (l *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(l.BasePath, path)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
