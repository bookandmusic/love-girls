package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/bookandmusic/love-girl/docs"
	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/handler"
	"github.com/bookandmusic/love-girl/internal/log"
	"github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
)

func ProvideGinEngine(
	cfg *config.AppConfig,
	logger *log.Logger,
) *server.GinEngine {
	return server.NewGinEngine(cfg, logger)
}

func ProvideRouter(
	cfg *config.AppConfig,
	ginEngine *server.GinEngine,
	authMiddleware *middleware.AuthMiddleware,
	handlers []handler.ApiHandler,
	staticHandlers []handler.StaticHandlerAware,
) *gin.Engine {
	engine := ginEngine.Engine

	// 初始化Swagger文档
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = cfg.Server.HostName

	// 循环注册静态路由（首页、Swagger UI 等）
	for _, h := range staticHandlers {
		h.RegisterRoutes(ginEngine)
	}

	apiGroup := engine.Group("/api/v1")
	// 循环注册 API 路由
	for _, h := range handlers {
		h.RegisterRoutes(apiGroup, ginEngine, authMiddleware)
	}

	return engine
}

func ProvideApp(
	cfg *config.AppConfig,
	logger *log.Logger,
	engine *gin.Engine,
	migrateErr error, // 添加 migrateErr 依赖
) *server.App {
	// 如果迁移失败，我们可以在这里处理错误
	if migrateErr != nil {
		logger.Error("Migration failed:", "error", migrateErr)
		// 这里可以选择 panic 或者以其他方式处理迁移错误
		// 但在依赖注入上下文中，我们通常只是返回应用实例
	}

	return server.NewApp(logger, engine, *cfg)
}

var RouterSet = wire.NewSet(
	ProvideGinEngine,
	ProvideRouter,
	ProvideApp,
)
