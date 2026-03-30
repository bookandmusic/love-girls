package provider

import (
	"github.com/google/wire"

	"github.com/bookandmusic/love-girl/internal/auth"
	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/repo"
	"github.com/bookandmusic/love-girl/internal/service"
	"github.com/bookandmusic/love-girl/internal/storage"
)

func ProvideUserService(log *log.Logger, userRepo *repo.UserRepo, fileRepo *repo.FileRepo, fileService *service.FileService, storage storage.Storage, cfg *config.AppConfig, jwt auth.JWT) *service.UserService {
	return service.NewUserService(log, *userRepo, *fileRepo, fileService, storage, &cfg.Server, jwt)
}

func ProvideFileService(log *log.Logger, storage storage.Storage, fileRepo *repo.FileRepo, cfg *config.AppConfig) *service.FileService {
	return service.NewFileService(log, storage, *fileRepo, &cfg.Server, &cfg.Storage, &cfg.ImageProxy)
}

func ProvideSystemService(
	log *log.Logger,
	userRepo *repo.UserRepo,
	settingRepo *repo.SettingRepo,
	albumRepo *repo.AlbumRepo,
	placeRepo *repo.PlaceRepo,
	momentRepo *repo.MomentRepo,
	fileService *service.FileService,
	cfg *config.AppConfig,
	jwt auth.JWT,
) *service.SystemService {
	return service.NewSystemService(log, *userRepo, *settingRepo, *albumRepo, *placeRepo, *momentRepo, fileService, cfg, jwt)
}

func ProvideAnniversaryService(log *log.Logger, anniversaryRepo *repo.AnniversaryRepo) *service.AnniversaryService {
	return service.NewAnniversaryService(log, anniversaryRepo)
}

func ProvideMomentService(log *log.Logger, momentRepo *repo.MomentRepo, commentRepo *repo.CommentRepo, fileService *service.FileService) *service.MomentService {
	return service.NewMomentService(log, momentRepo, commentRepo, fileService)
}

func ProvidePlaceService(log *log.Logger, placeRepo *repo.PlaceRepo, fileService *service.FileService) *service.PlaceService {
	return service.NewPlaceService(log, placeRepo, fileService)
}

func ProvideAlbumService(log *log.Logger, albumRepo *repo.AlbumRepo, fileService *service.FileService) *service.AlbumService {
	return service.NewAlbumService(log, albumRepo, fileService)
}

func ProvideCommentService(log *log.Logger, commentRepo *repo.CommentRepo, momentRepo *repo.MomentRepo, notificationRepo *repo.NotificationRepo, fileService *service.FileService, notificationService *service.NotificationService) *service.CommentService {
	return service.NewCommentService(log, commentRepo, momentRepo, notificationRepo, fileService, notificationService)
}

func ProvideNotificationService(log *log.Logger, notificationRepo *repo.NotificationRepo, fileService *service.FileService) *service.NotificationService {
	return service.NewNotificationService(log, notificationRepo, fileService)
}

var ServiceSet = wire.NewSet(
	ProvideUserService,
	ProvideFileService,
	ProvideSystemService,
	ProvideAnniversaryService,
	ProvideMomentService,
	ProvidePlaceService,
	ProvideAlbumService,
	ProvideCommentService,
	ProvideNotificationService,
)
