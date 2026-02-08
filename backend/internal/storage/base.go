package storage

import (
	"context"
	"io"
)

type GinProxyURLBuilder func(fileID uint64) string

type Storage interface {
	Name() string
	Save(ctx context.Context, path string, r io.Reader) error
	Open(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
	URL(ctx context.Context, fileID uint64, filePath string, width, height int, builder GinProxyURLBuilder) (string, error)
}
