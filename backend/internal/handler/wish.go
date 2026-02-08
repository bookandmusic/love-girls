package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/love-girl/internal/auth"
	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/internal/service"
)

type WishHandler struct {
	WishService *service.WishService
	JWT         auth.JWT
}

func NewWishHandler(wishService *service.WishService, jwt auth.JWT) *WishHandler {
	return &WishHandler{
		WishService: wishService,
		JWT:         jwt,
	}
}

// RegisterRoutes 注册祝福相关的路由
func (h *WishHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	wishGroup := apiGroup.Group("/wishes")
	{
		// 不需要认证的路由
		wishGroup.GET("/", h.ListWishs)   // 获取祝福列表
		wishGroup.POST("/", h.CreateWish) // 创建祝福

		// 需要认证的路由
		authGroup := wishGroup.Group("")
		authGroup.Use(authMiddleware.Handle())
		{
			authGroup.DELETE("/:id", h.DeleteWish)       // 删除祝福
			authGroup.PUT("/:id/approve", h.ApproveWish) // 批准祝福
		}
	}
}

// ListWishs 获取祝福列表
// @Summary 获取祝福列表
// @Description 获取祝福列表，支持分页查询。未认证用户只能查看已审核的祝福，已认证用户可以查看所有祝福。
// @Tags wishes
// @Accept json
// @Produce json
// @Param page query int false "页码，从 1 开始"
// @Param size query int false "每页数量"
// @Param approved query bool false "是否只获取已审核的祝福，仅对已认证用户有效"
// @Success 200 {object} Response{data=service.WishListResponse}
// @Router /api/v1/wishes [get]
func (h *WishHandler) ListWishs(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// 尝试从 Authorization header 中获取用户认证状态（可选认证）
	var isAuthenticated bool
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			token := strings.TrimSpace(parts[1])
			_, err := h.JWT.Parse(token)
			if err == nil {
				isAuthenticated = true
			}
		}
	}

	var approved *bool
	approvedStr := c.Query("approved")

	if isAuthenticated {
		// 已认证用户：根据参数决定
		if approvedStr != "" {
			approvedValue, err := strconv.ParseBool(approvedStr)
			if err == nil {
				approved = &approvedValue
			}
		}
		// 如果没有参数，approved为nil，返回所有祝福
	} else {
		// 未认证用户：只能查看已审核的祝福
		defaultApproved := true
		approved = &defaultApproved
	}

	// 调用服务层获取祝福列表
	wish, err := h.WishService.ListWishs(ctx, page, size, approved)
	if err != nil {
		h.WishService.Log.Error("获取祝福列表失败", "error", err)
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
		Data:    wish,
	})
}

// CreateWish 创建祝福
// @Summary 创建祝福
// @Description 创建一个新的祝福
// @Tags wishes
// @Accept json
// @Produce json
// @Param blessing body service.WishCreateRequest true "祝福信息"
// @Success 200 {object} Response{data=service.Wish}
// @Router /api/v1/wishes [post]
func (h *WishHandler) CreateWish(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析请求参数
	var req service.WishCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.WishService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 调用服务层创建祝福
	blessing, err := h.WishService.CreateWish(ctx, &req)
	if err != nil {
		h.WishService.Log.Error("创建祝福失败", "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "创建成功",
		Data:    blessing,
	})
}

// DeleteWish 删除祝福
// @Summary 删除祝福
// @Description 删除指定的祝福
// @Tags wishes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "祝福 ID"
// @Success 200 {object} Response{data=nil}
// @Router /api/v1/wishes/{id} [delete]
func (h *WishHandler) DeleteWish(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析祝福ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.WishService.Log.Error("无效的祝福ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的祝福ID",
			Data:    nil,
		})
		return
	}

	// 调用服务层删除祝福
	ok, err := h.WishService.DeleteWish(ctx, id)
	if err != nil {
		h.WishService.Log.Error("删除祝福失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if !ok {
		h.WishService.Log.Info("祝福不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "祝福不存在",
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

// ApproveWish 批准祝福
// @Summary 批准祝福
// @Description 批准指定的祝福，使其可以公开显示
// @Tags wishes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "祝福 ID"
// @Success 200 {object} Response{data=nil}
// @Router /api/v1/wishes/{id}/approve [put]
func (h *WishHandler) ApproveWish(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析祝福ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.WishService.Log.Error("无效的祝福ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的祝福ID",
			Data:    nil,
		})
		return
	}

	// 调用服务层批准祝福
	err = h.WishService.ApproveWish(ctx, id)
	if err != nil {
		h.WishService.Log.Error("批准祝福失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "审核成功",
		Data:    nil,
	})
}
