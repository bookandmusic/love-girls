package provider

import (
	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/storage"
)

func ProvideStorageLocal(
	cfg *config.AppConfig,
	logger *log.Logger,
) (*storage.LocalStorage, error) {
	// Local 存储路径由 data_dir 自动计算
	uploadDir := cfg.GetDataPaths().UploadDir
	return storage.NewLocalStorage(uploadDir)
}

func ProvideStorageS3(
	cfg *config.AppConfig,
	logger *log.Logger,
) (*storage.S3Storage, error) {
	return storage.NewS3Storage(cfg.Storage.S3)
}

func ProvideStorageWebDAV(
	cfg *config.AppConfig,
	logger *log.Logger,
) (*storage.WebDAVStorage, error) {
	return storage.NewWebDAVStorage(cfg.Storage.WebDAV)
}

func ProvideStorage(
	cfg *config.AppConfig,
	logger *log.Logger,
) (storage.Storage, error) {
	switch cfg.Storage.Backend {
	case "local":
		return ProvideStorageLocal(cfg, logger)
	case "s3":
		return ProvideStorageS3(cfg, logger)
	case "webdav":
		return ProvideStorageWebDAV(cfg, logger)
	default:
		return ProvideStorageLocal(cfg, logger)
	}
}
