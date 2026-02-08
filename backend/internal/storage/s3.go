package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/bookandmusic/love-girl/internal/config"
)

type S3Storage struct {
	cfg    *config.S3StorageConfig
	access *config.StorageAccessConfig
	client *minio.Client
}

func NewS3Storage(cfg *config.S3StorageConfig, access *config.StorageAccessConfig) (*S3Storage, error) {

	// 创建 MinIO 客户端
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Credentials.AccessKeyID, cfg.Credentials.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// 如果是自定义 endpoint（如 MinIO），建议启用 bucket lookup via path style
	// MinIO 默认支持，无需额外设置（minio-go 默认就是 path-style）

	return &S3Storage{
		cfg:    cfg,
		access: access,
		client: client,
	}, nil
}

func (s *S3Storage) Name() string {
	return "s3" // 或仍返回 "s3"，按需
}

// Save 上传文件
func (s *S3Storage) Save(ctx context.Context, path string, r io.Reader) error {
	// 获取 reader 长度（可选优化）；若无法获取，MinIO 会自动分块上传
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
func (s *S3Storage) presignClient() (*minio.Client, error) {
	// 如果没有 PublicBaseURL，直接复用当前 client
	if s.cfg.PublicBaseURL == "" {
		return s.client, nil
	}

	u, err := url.Parse(s.cfg.PublicBaseURL)
	if err != nil {
		return nil, err
	}

	endpoint := u.Host
	useSSL := u.Scheme == "https"

	return minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			s.cfg.Credentials.AccessKeyID,
			s.cfg.Credentials.SecretAccessKey,
			"",
		),
		Secure: useSSL,
		Region: s.cfg.Region,
	})
}

func (s *S3Storage) URL(
	ctx context.Context,
	fileID uint64, filePath string,
	width, height int,
	builder GinProxyURLBuilder,
) (string, error) {

	var imageURL string

	// 1. Gin Proxy 优先
	if s.access.GinProxy.Enabled {
		imageURL = builder(fileID)
	} else {
		// 2. Presign
		client, err := s.presignClient()
		if err != nil {
			return "", err
		}

		expire := time.Duration(s.cfg.PresignExpire) * time.Second
		if expire <= 0 {
			expire = 7 * 24 * time.Hour
		}

		presignedURL, err := client.PresignedGetObject(
			ctx,
			s.cfg.Bucket,
			filePath,
			expire,
			nil,
		)
		if err != nil {
			return "", err
		}

		imageURL = presignedURL.String()
	}

	// 3. Image Proxy（可选）
	if s.access.ImageProxy.Enabled {
		escaped := url.QueryEscape(imageURL)
		return fmt.Sprintf(
			"%s/%dx%d/%s",
			s.access.ImageProxy.BaseURL,
			width,
			height,
			escaped,
		), nil
	}

	return imageURL, nil
}

// Delete 删除文件
func (s *S3Storage) Delete(ctx context.Context, path string) error {
	return s.client.RemoveObject(ctx, s.cfg.Bucket, path, minio.RemoveObjectOptions{})
}
