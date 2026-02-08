package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// RegisterRoutes 注册健康检查路由
func (h *HealthHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	apiGroup.GET("/health", h.HealthCheck)
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "OK",
	})
}
