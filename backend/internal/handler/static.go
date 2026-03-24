package handler

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

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

	// 获取两个前端的静态文件系统
	frontendFS, frontendErr := fs.Sub(h.fs, "assets/dist-frontend")
	adminFS, adminErr := fs.Sub(h.fs, "assets/dist-admin")

	// 读取两个前端的 index.html
	frontendIndex, frontendIndexErr := h.fs.ReadFile("assets/dist-frontend/index.html")
	adminIndex, adminIndexErr := h.fs.ReadFile("assets/dist-admin/index.html")

	// 如果两个前端都不存在，返回错误
	if frontendErr != nil || frontendIndexErr != nil {
		engine.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "frontend static files not found"})
		})
		return
	}

	frontendFileServer := http.FileServer(http.FS(frontendFS))
	adminFileServer := http.FileServer(http.FS(adminFS))

	// 首页路由 - 前台
	engine.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", frontendIndex)
	})

	// /admin 路由 - 后台
	engine.GET("/admin", func(c *gin.Context) {
		if adminErr != nil || adminIndexErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "admin frontend not found"})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", adminIndex)
	})

	// /admin/ 路由 - 后台
	engine.GET("/admin/", func(c *gin.Context) {
		if adminErr != nil || adminIndexErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "admin frontend not found"})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", adminIndex)
	})

	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 处理 /admin 路径下的请求
		if strings.HasPrefix(path, "/admin/") {
			relPath := strings.TrimPrefix(path, "/admin/")

			// 如果是空路径，返回 admin index.html
			if relPath == "" {
				c.Data(http.StatusOK, "text/html; charset=utf-8", adminIndex)
				return
			}

			// 尝试打开静态文件
			file, err := adminFS.Open(relPath)
			if err == nil {
				file.Close()
				c.Header("Cache-Control", "public, max-age=31536000")
				adminFileServer.ServeHTTP(c.Writer, c.Request)
				return
			}

			// 文件不存在，返回 admin index.html 让前端路由处理
			c.Data(http.StatusOK, "text/html; charset=utf-8", adminIndex)
			return
		}

		// 处理其他路径 - 前台
		relPath := strings.TrimPrefix(path, "/")

		// 尝试打开静态文件
		file, err := frontendFS.Open(relPath)
		if err == nil {
			file.Close()
			c.Header("Cache-Control", "public, max-age=31536000")
			frontendFileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// 文件不存在，返回前台 index.html 让前端路由处理
		c.Data(http.StatusOK, "text/html; charset=utf-8", frontendIndex)
	})
}

// SwaggerHandler Swagger UI 处理器
type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) RegisterRoutes(ginEngine *server.GinEngine) {
	engine := ginEngine.Engine
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
