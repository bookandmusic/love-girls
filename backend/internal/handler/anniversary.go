package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/internal/service"
)

type AnniversaryHandler struct {
	AnniversaryService *service.AnniversaryService
}

func NewAnniversaryHandler(anniversaryService *service.AnniversaryService) *AnniversaryHandler {
	return &AnniversaryHandler{
		AnniversaryService: anniversaryService,
	}
}

// RegisterRoutes 注册纪念日相关的路由
func (h *AnniversaryHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	anniversaryGroup := apiGroup.Group("/anniversaries")
	{
		// 不需要认证的路由
		anniversaryGroup.GET("/", h.ListAnniversaries) // 获取纪念日列表

		// 需要认证的路由
		authGroup := anniversaryGroup.Group("")
		authGroup.Use(authMiddleware.Handle())
		{
			authGroup.POST("/", h.CreateAnniversary)      // 创建纪念日
			authGroup.PUT("/:id", h.UpdateAnniversary)    // 更新纪念日
			authGroup.DELETE("/:id", h.DeleteAnniversary) // 删除纪念日
		}
	}
}

// ListAnniversaries 获取纪念日列表
// @Summary 获取纪念日列表
// @Description 分页获取纪念日列表
// @Tags anniversaries
// @Produce json
// @Param page query int true "页码" default(1)
// @Param size query int true "每页数量" default(10)
// @Success 200 {object} Response{data=service.AnniversaryListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /anniversaries [get]
func (h *AnniversaryHandler) ListAnniversaries(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "页码参数无效",
			Data:    nil,
		})
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil || size < 1 || size > 100 {
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "每页数量参数无效",
			Data:    nil,
		})
		return
	}

	// 调用服务层获取纪念日列表
	response, err := h.AnniversaryService.ListAnniversaries(ctx, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "获取纪念日列表失败",
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

// CreateAnniversary 创建纪念日
// @Summary 创建纪念日
// @Description 创建一个新的纪念日
// @Tags anniversaries
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param anniversary body service.AnniversaryCreateRequest true "纪念日信息"
// @Success 200 {object} Response{data=service.FrontendAnniversary}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /anniversaries [post]
func (h *AnniversaryHandler) CreateAnniversary(c *gin.Context) {
	ctx := c.Request.Context()

	var req service.AnniversaryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.AnniversaryService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "参数格式错误或字段缺失",
			Data:    nil,
		})
		return
	}

	// 调用服务层创建纪念日
	anniversary, err := h.AnniversaryService.CreateAnniversary(ctx, &req)
	if err != nil {
		h.AnniversaryService.Log.Error("创建纪念日失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "创建纪念日失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "创建成功",
		Data:    anniversary,
	})
}

// UpdateAnniversary 更新纪念日
// @Summary 更新纪念日
// @Description 更新指定ID的纪念日信息
// @Tags anniversaries
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path int true "纪念日ID"
// @Param anniversary body service.AnniversaryUpdateRequest true "纪念日更新信息"
// @Success 200 {object} Response{data=service.FrontendAnniversary}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /anniversaries/{id} [put]
func (h *AnniversaryHandler) UpdateAnniversary(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "纪念日ID无效",
			Data:    nil,
		})
		return
	}

	var req service.AnniversaryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.AnniversaryService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "参数格式错误或字段缺失",
			Data:    nil,
		})
		return
	}

	// 调用服务层更新纪念日
	anniversary, err := h.AnniversaryService.UpdateAnniversary(ctx, id, &req)
	if err != nil {
		if err.Error() == "纪念日不存在" {
			c.JSON(http.StatusNotFound, Response{
				Code:    1,
				Message: "纪念日不存在",
				Data:    nil,
			})
			return
		}

		h.AnniversaryService.Log.Error("更新纪念日失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "更新纪念日失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "更新成功",
		Data:    anniversary,
	})
}

// DeleteAnniversary 删除纪念日
// @Summary 删除纪念日
// @Description 删除指定ID的纪念日
// @Tags anniversaries
// @Produce json
// @Security OAuth2Password
// @Param id path int true "纪念日ID"
// @Success 200 {object} Response{data=nil}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /anniversaries/{id} [delete]
func (h *AnniversaryHandler) DeleteAnniversary(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "纪念日ID无效",
			Data:    nil,
		})
		return
	}

	// 调用服务层删除纪念日
	err = h.AnniversaryService.DeleteAnniversary(ctx, id)
	if err != nil {
		if err.Error() == "纪念日不存在" {
			c.JSON(http.StatusNotFound, Response{
				Code:    1,
				Message: "纪念日不存在",
				Data:    nil,
			})
			return
		}

		h.AnniversaryService.Log.Error("删除纪念日失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "删除纪念日失败",
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
