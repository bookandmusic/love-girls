package server

import (
	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/log"
)

// GinEngine 封装 Gin 引擎和日志
type GinEngine struct {
	Engine *gin.Engine
	Config *config.AppConfig
	Logger *log.Logger
}

// NewGinEngine 构造函数，初始化 Gin 引擎和配置
func NewGinEngine(cfg *config.AppConfig, logger *log.Logger) *GinEngine {

	// 初始化 Gin 引擎
	engine := gin.New()

	// 设置文件上传大小限制为200MB
	engine.MaxMultipartMemory = 200 << 20 // 200 MB

	// 使用自定义的日志中间件
	engine.Use(gin.LoggerWithWriter(logger)) // 通过自定义的 Logger 输出日志
	engine.Use(gin.Recovery())

	// 创建 GinEngine 实例
	return &GinEngine{
		Engine: engine,
		Config: cfg,
		Logger: logger,
	}
}
