package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/bookandmusic/love-girl/internal/config"
)

type LocalStorage struct {
	Root   string
	access *config.StorageAccessConfig
}

func NewLocalStorage(root string, access *config.StorageAccessConfig) (*LocalStorage, error) {
	// 初始化时创建根目录
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("创建存储根目录失败: %w", err)
	}

	return &LocalStorage{
		Root:   root,
		access: access,
	}, nil
}

func (l *LocalStorage) Name() string {
	return "local"
}

func (l *LocalStorage) Save(ctx context.Context, filePath string, r io.Reader) error {
	fullPath := filepath.Join(l.Root, filePath)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}

func (l *LocalStorage) Open(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.Root, filePath)
	return os.Open(fullPath)
}

func (l *LocalStorage) URL(ctx context.Context, fileID uint64, filePath string, width, height int, builder GinProxyURLBuilder) (string, error) {
	imageUrl := builder(fileID)
	if l.access.ImageProxy.Enabled {
		imageUrl = url.QueryEscape(imageUrl)
		return fmt.Sprintf("%s/%dx%d/%s", l.access.ImageProxy.BaseURL, width, height, imageUrl), nil
	}
	return imageUrl, nil
}

func (l *LocalStorage) Delete(ctx context.Context, filePath string) error {
	fullPath := filepath.Join(l.Root, filePath)
	return os.Remove(fullPath)
}
