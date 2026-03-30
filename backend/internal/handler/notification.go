package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/love-girl/internal/auth"
	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/internal/service"
)

type NotificationHandler struct {
	NotificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		NotificationService: notificationService,
	}
}

func (h *NotificationHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	authGroup := apiGroup.Group("")
	authGroup.Use(authMiddleware.Handle())
	{
		authGroup.GET("/notifications/unread", h.ListUnreadNotifications)
		authGroup.POST("/notifications/:id/read", h.MarkAsRead)
		authGroup.GET("/notifications/count", h.GetUnreadCount)
		authGroup.POST("/notifications/read-all", h.MarkAllAsRead)
	}
}

func (h *NotificationHandler) ListUnreadNotifications(c *gin.Context) {
	claims, ok := auth.GetAuthClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    1,
			Message: "未登录",
			Data:    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	response, err := h.NotificationService.ListUnreadNotifications(c, claims.UserID, page, size)
	if err != nil {
		h.NotificationService.Log.Error("获取通知列表失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "查询成功",
		Data:    response,
	})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	ctx := c.Request.Context()

	claims, ok := auth.GetAuthClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    1,
			Message: "未登录",
			Data:    nil,
		})
		return
	}

	count, err := h.NotificationService.GetUnreadCount(ctx, claims.UserID)
	if err != nil {
		h.NotificationService.Log.Error("获取未读数量失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "查询成功",
		Data: map[string]int64{
			"count": count,
		},
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.NotificationService.Log.Error("无效的通知ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的通知ID",
			Data:    nil,
		})
		return
	}

	claims, ok := auth.GetAuthClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    1,
			Message: "未登录",
			Data:    nil,
		})
		return
	}

	if err := h.NotificationService.MarkAsRead(ctx, id, claims.UserID); err != nil {
		h.NotificationService.Log.Error("标记已读失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "标记成功",
		Data:    nil,
	})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	ctx := c.Request.Context()

	claims, ok := auth.GetAuthClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    1,
			Message: "未登录",
			Data:    nil,
		})
		return
	}

	if err := h.NotificationService.MarkAllAsRead(ctx, claims.UserID); err != nil {
		h.NotificationService.Log.Error("全部标记已读失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "标记成功",
		Data:    nil,
	})
}
