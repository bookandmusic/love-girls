package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/internal/service"
)

type PlaceHandler struct {
	PlaceService *service.PlaceService
}

func NewPlaceHandler(placeService *service.PlaceService) *PlaceHandler {
	return &PlaceHandler{
		PlaceService: placeService,
	}
}

// RegisterRoutes 注册地点相关的路由
func (h *PlaceHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	placeGroup := apiGroup.Group("/places")
	{
		// 不需要认证的路由
		placeGroup.GET("/", h.ListPlaces) // 获取所有地点

		// 需要认证的路由
		authGroup := placeGroup.Group("")
		authGroup.Use(authMiddleware.Handle())
		{
			authGroup.POST("/", h.CreatePlace)      // 创建地点
			authGroup.PUT("/:id", h.UpdatePlace)    // 更新地点
			authGroup.DELETE("/:id", h.DeletePlace) // 删除地点
		}
	}
}

// ListPlaces 获取所有地点
// @Summary 获取所有地点
// @Description 获取所有地点，保持分页数据结构
// @Tags places
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=service.PlaceListResponse}
// @Router /api/v1/places [get]
func (h *PlaceHandler) ListPlaces(c *gin.Context) {
	ctx := c.Request.Context()

	// 调用服务层获取所有地点
	places, err := h.PlaceService.ListPlaces(ctx)
	if err != nil {
		h.PlaceService.Log.Error("获取地点列表失败", "error", err)
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
		Data:    places,
	})
}

// CreatePlace 创建地点
// @Summary 创建地点
// @Description 创建一个新的地点
// @Tags places
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param place body service.PlaceCreateRequest true "地点信息"
// @Success 200 {object} Response{data=service.PlaceResponse}
// @Router /api/v1/places [post]
func (h *PlaceHandler) CreatePlace(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析请求参数
	var req service.PlaceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.PlaceService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 调用服务层创建地点
	place, err := h.PlaceService.CreatePlace(ctx, &req)
	if err != nil {
		h.PlaceService.Log.Error("创建地点失败", "error", err, "request", req)
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
		Data:    place,
	})
}

// UpdatePlace 更新地点
// @Summary 更新地点
// @Description 更新地点的信息
// @Tags places
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "地点ID"
// @Param place body service.PlaceUpdateRequest true "地点信息"
// @Success 200 {object} Response{data=service.PlaceResponse}
// @Router /api/v1/places/{id} [put]
func (h *PlaceHandler) UpdatePlace(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析地点ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.PlaceService.Log.Error("无效的地点ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的地点ID",
			Data:    nil,
		})
		return
	}

	// 解析请求参数
	var req service.PlaceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.PlaceService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 调用服务层更新地点
	place, err := h.PlaceService.UpdatePlace(ctx, id, &req)
	if err != nil {
		h.PlaceService.Log.Error("更新地点失败", "id", id, "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if place == nil {
		h.PlaceService.Log.Info("地点不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "地点不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "更新成功",
		Data:    place,
	})
}

// DeletePlace 删除地点
// @Summary 删除地点
// @Description 删除指定的地点
// @Tags places
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "地点ID"
// @Success 200 {object} Response{data=nil}
// @Router /api/v1/places/{id} [delete]
func (h *PlaceHandler) DeletePlace(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析地点ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.PlaceService.Log.Error("无效的地点ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的地点ID",
			Data:    nil,
		})
		return
	}

	// 调用服务层删除地点
	ok, err := h.PlaceService.DeletePlace(ctx, id)
	if err != nil {
		h.PlaceService.Log.Error("删除地点失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if !ok {
		h.PlaceService.Log.Info("地点不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "地点不存在",
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
