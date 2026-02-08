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
	return storage.NewLocalStorage(cfg.Storage.Local.Root, &cfg.Storage.Access)
}

func ProvideStorageS3(
	cfg *config.AppConfig,
	logger *log.Logger,
) (*storage.S3Storage, error) {
	return storage.NewS3Storage(cfg.Storage.S3, &cfg.Storage.Access)
}

func ProvideStorageWebDAV(
	cfg *config.AppConfig,
	logger *log.Logger,
) (*storage.WebDAVStorage, error) {
	return storage.NewWebDAVStorage(cfg.Storage.WebDAV, &cfg.Storage.Access)
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
