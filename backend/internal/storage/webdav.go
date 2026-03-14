package storage

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/studio-b12/gowebdav"

	"github.com/bookandmusic/love-girl/internal/config"
)

type WebDAVStorage struct {
	cfg    *config.WebDAVStorageConfig
	client *gowebdav.Client
}

func NewWebDAVStorage(cfg *config.WebDAVStorageConfig) (*WebDAVStorage, error) {
	client := gowebdav.NewClient(cfg.Endpoint, cfg.Auth.Username, cfg.Auth.Password)
	// 测试连接
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to create webdav client: %w", err)
	}

	return &WebDAVStorage{
		cfg:    cfg,
		client: client,
	}, nil
}

func (w *WebDAVStorage) Name() string {
	return "webdav"
}

func (w *WebDAVStorage) Save(ctx context.Context, filePath string, r io.Reader) error {
	fullPath := path.Join(w.cfg.BasePath, filePath)
	return w.client.WriteStream(fullPath, r, 0644)
}

func (w *WebDAVStorage) Open(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fullPath := path.Join(w.cfg.BasePath, filePath)
	return w.client.ReadStream(fullPath)
}

func (w *WebDAVStorage) URL(ctx context.Context, fileID uint64, filePath string, width, height int, builder GinProxyURLBuilder) (string, error) {
	if w.cfg.PublicURL != "" {
		return fmt.Sprintf("%s/%s/%s", w.cfg.PublicURL, w.cfg.BasePath, filePath), nil
	}

	return builder(fileID), nil
}

func (w *WebDAVStorage) Delete(ctx context.Context, filePath string) error {
	fullPath := path.Join(w.cfg.BasePath, filePath)
	return w.client.Remove(fullPath)
}