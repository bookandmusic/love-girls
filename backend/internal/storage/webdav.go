package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"

	"github.com/studio-b12/gowebdav"

	"github.com/bookandmusic/love-girl/internal/config"
)

type WebDAVStorage struct {
	cfg    *config.WebDAVStorageConfig
	access *config.StorageAccessConfig
	client *gowebdav.Client
}

func NewWebDAVStorage(cfg *config.WebDAVStorageConfig, access *config.StorageAccessConfig) (*WebDAVStorage, error) {
	client := gowebdav.NewClient(cfg.Endpoint, cfg.Auth.Username, cfg.Auth.Password)
	// 测试连接
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to create webdav client: %w", err)
	}

	return &WebDAVStorage{
		cfg:    cfg,
		access: access,
		client: client,
	}, nil
}

func (w *WebDAVStorage) Name() string {
	return "webdav"
}

func (w *WebDAVStorage) Save(ctx context.Context, filePath string, r io.Reader) error {
	fullPath := path.Join(w.cfg.BasePath, filePath)
	// 上传时直接覆盖
	return w.client.WriteStream(fullPath, r, 0644)
}

func (w *WebDAVStorage) Open(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fullPath := path.Join(w.cfg.BasePath, filePath)
	return w.client.ReadStream(fullPath)
}

func (w *WebDAVStorage) URL(ctx context.Context, fileID uint64, filePath string, width, height int, builder GinProxyURLBuilder) (string, error) {
	var imageUrl string
	if w.access.GinProxy.Enabled {
		imageUrl = builder(fileID)
	} else {
		u, err := url.Parse(w.cfg.PublicBaseURL)
		if err != nil {
			return "", err
		}
		if w.cfg.Auth.Username != "" && w.cfg.Auth.Password != "" {
			u.User = url.UserPassword(w.cfg.Auth.Username, w.cfg.Auth.Password)
		}
		imageUrl, err = url.JoinPath(u.String(), w.cfg.BasePath, filePath)
		if err != nil {
			return "", err
		}
	}
	if w.access.ImageProxy.Enabled {
		imageUrl = url.QueryEscape(imageUrl)
		return fmt.Sprintf("%s/%dx%d/%s", w.access.ImageProxy.BaseURL, width, height, imageUrl), nil
	}
	return imageUrl, nil
}

func (w *WebDAVStorage) Delete(ctx context.Context, filePath string) error {
	fullPath := path.Join(w.cfg.BasePath, filePath)
	return w.client.Remove(fullPath)
}
