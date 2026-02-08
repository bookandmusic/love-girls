package handler

import (
	"github.com/gin-gonic/gin"

	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
)

// ApiHandler 定义 API 路由处理器接口
type ApiHandler interface {
	RegisterRoutes(engine *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware)
}

// StaticHandlerAware 静态文件处理器接口（直接在 Gin 引擎根路径注册）
type StaticHandlerAware interface {
	RegisterRoutes(ginEngine *server.GinEngine)
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
