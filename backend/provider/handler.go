package provider

import (
	"github.com/google/wire"

	"github.com/bookandmusic/love-girl/internal/auth"
	"github.com/bookandmusic/love-girl/internal/handler"
	"github.com/bookandmusic/love-girl/internal/service"
)

func ProvideUserHandler(svc *service.UserService) *handler.UserHandler {
	return handler.NewUserHandler(svc)
}

func ProvideHealthHandler() *handler.HealthHandler {
	return handler.NewHealthHandler()
}

func ProvideFileHandler(svc *service.FileService) *handler.FileHandler {
	return handler.NewFileHandler(svc)
}

func ProvideSystemHandler(svc *service.SystemService) *handler.SystemHandler {
	return handler.NewSystemHandler(svc)
}

func ProvideMomentHandler(svc *service.MomentService) *handler.MomentHandler {
	return handler.NewMomentHandler(svc)
}

func ProvideAnniversaryHandler(svc *service.AnniversaryService) *handler.AnniversaryHandler {
	return handler.NewAnniversaryHandler(svc)
}

func ProvidePlaceHandler(svc *service.PlaceService) *handler.PlaceHandler {
	return handler.NewPlaceHandler(svc)
}

func ProvideWishHandler(svc *service.WishService, jwt auth.JWT) *handler.WishHandler {
	return handler.NewWishHandler(svc, jwt)
}

func ProvideAlbumHandler(svc *service.AlbumService) *handler.AlbumHandler {
	return handler.NewAlbumHandler(svc)
}

func ProvideStaticHandler() *handler.StaticHandler {
	return handler.NewStaticHandler()
}

func ProvideSwaggerHandler() *handler.SwaggerHandler {
	return handler.NewSwaggerHandler()
}

func ProvideStaticHandlers(
	staticHandler *handler.StaticHandler,
	swaggerHandler *handler.SwaggerHandler,
) []handler.StaticHandlerAware {
	return []handler.StaticHandlerAware{
		staticHandler,
		swaggerHandler,
	}
}

func ProvideHandlers(
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	fileHandler *handler.FileHandler,
	systemHandler *handler.SystemHandler,
	momentHandler *handler.MomentHandler,
	anniversaryHandler *handler.AnniversaryHandler,
	placeHandler *handler.PlaceHandler,
	wishHandler *handler.WishHandler,
	albumHandler *handler.AlbumHandler,
) []handler.ApiHandler {
	return []handler.ApiHandler{
		userHandler,
		healthHandler,
		fileHandler,
		systemHandler,
		momentHandler,
		anniversaryHandler,
		placeHandler,
		wishHandler,
		albumHandler,
	}
}

var HandlerSet = wire.NewSet(
	ProvideUserHandler,
	ProvideHealthHandler,
	ProvideFileHandler,
	ProvideSystemHandler,
	ProvideMomentHandler,
	ProvideAnniversaryHandler,
	ProvidePlaceHandler,
	ProvideWishHandler,
	ProvideAlbumHandler,
	ProvideStaticHandler,
	ProvideSwaggerHandler,
	ProvideStaticHandlers,
	ProvideHandlers,
)
