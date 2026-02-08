package handler

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/bookandmusic/love-girl/internal/server"
)

type StaticHandler struct {
	fs embed.FS
}

func NewStaticHandler() *StaticHandler {
	return &StaticHandler{fs: distFS}
}

func (h *StaticHandler) RegisterRoutes(ginEngine *server.GinEngine) {
	engine := ginEngine.Engine

	subFS, err := fs.Sub(h.fs, "assets/dist")
	if err != nil {
		engine.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "static files not found"})
		})
		return
	}

	indexData, err := h.fs.ReadFile("assets/dist/index.html")
	if err != nil {
		engine.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "index.html not found"})
		})
		return
	}

	engine.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexData)
	})

	fileServer := http.FileServer(http.FS(subFS))

	engine.NoRoute(func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=31536000")
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

// SwaggerHandler Swagger UI 处理器
type SwaggerHandler struct {
}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) RegisterRoutes(ginEngine *server.GinEngine) {
	engine := ginEngine.Engine
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
