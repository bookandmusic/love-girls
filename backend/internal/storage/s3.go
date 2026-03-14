package storage

import (
	"context"
	"fmt"
	"io"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/bookandmusic/love-girl/internal/config"
)

type S3Storage struct {
	cfg    *config.S3StorageConfig
	client *minio.Client
}

func NewS3Storage(cfg *config.S3StorageConfig) (*S3Storage, error) {
	// 创建 MinIO 客户端
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Credentials.AccessKeyID, cfg.Credentials.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &S3Storage{
		cfg:    cfg,
		client: client,
	}, nil
}

func (s *S3Storage) Name() string {
	return "s3"
}

// Save 上传文件
func (s *S3Storage) Save(ctx context.Context, path string, r io.Reader) error {
	_, err := s.client.PutObject(ctx, s.cfg.Bucket, path, r, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

// Open 以代理方式读取文件（返回 io.ReadCloser）
func (s *S3Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.cfg.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	// 检查对象是否存在（GetObject 不会立即报错，需调用 Stat）
	if _, err = obj.Stat(); err != nil {
		obj.Close()
		return nil, err
	}

	return obj, nil
}

// URL 返回访问 URL
func (s *S3Storage) URL(
	ctx context.Context,
	fileID uint64, filePath string,
	width, height int,
	builder GinProxyURLBuilder,
) (string, error) {
	if s.cfg.PublicURL != "" {
		return fmt.Sprintf("%s/%s", s.cfg.PublicURL, filePath), nil
	}

	return builder(fileID), nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	return s.client.RemoveObject(ctx, s.cfg.Bucket, path, minio.RemoveObjectOptions{})
}
