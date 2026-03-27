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

type CommentHandler struct {
	CommentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
	}
}

func (h *CommentHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	apiGroup.GET("/moments/:id/comments", h.ListComments)

	authGroup := apiGroup.Group("")
	authGroup.Use(authMiddleware.Handle())
	{
		authGroup.POST("/moments/:id/comments", h.CreateComment)
		authGroup.DELETE("/comments/:id", h.DeleteComment)
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	momentIDStr := c.Param("id")
	momentID, err := strconv.ParseUint(momentIDStr, 10, 64)
	if err != nil {
		h.CommentService.Log.Error("无效的动态ID", "id", momentIDStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的动态ID",
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

	var req service.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.CommentService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "参数校验失败",
			Data:    nil,
		})
		return
	}

	req.MomentID = momentID
	req.UserID = claims.UserID

	comment, err := h.CommentService.CreateComment(c, &req)
	if err != nil {
		h.CommentService.Log.Error("创建评论失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "评论成功",
		Data:    comment,
	})
}

func (h *CommentHandler) ListComments(c *gin.Context) {
	momentIDStr := c.Param("id")
	momentID, err := strconv.ParseUint(momentIDStr, 10, 64)
	if err != nil {
		h.CommentService.Log.Error("无效的动态ID", "id", momentIDStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的动态ID",
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

	response, err := h.CommentService.ListComments(c, momentID, page, size)
	if err != nil {
		h.CommentService.Log.Error("获取评论列表失败", "error", err, "momentID", momentID)
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

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.CommentService.Log.Error("无效的评论ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的评论ID",
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

	deleted, err := h.CommentService.DeleteComment(ctx, id, claims.UserID)
	if err != nil {
		h.CommentService.Log.Error("删除评论失败", "id", id, "error", err)
		if err.Error() == "无权限删除此评论" {
			c.JSON(http.StatusForbidden, Response{
				Code:    1,
				Message: "无权限删除此评论",
				Data:    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if !deleted {
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "评论不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "删除成功",
		Data:    nil,
	})
}
