package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/internal/service"
)

type AlbumHandler struct {
	AlbumService *service.AlbumService
}

func NewAlbumHandler(albumService *service.AlbumService) *AlbumHandler {
	return &AlbumHandler{
		AlbumService: albumService,
	}
}

// RegisterRoutes 注册相册相关的路由
func (h *AlbumHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	albumGroup := apiGroup.Group("/albums")
	{
		// 不需要认证的路由
		albumGroup.GET("/", h.ListAlbums)                // 获取相册列表
		albumGroup.GET("/:id/photos", h.ListAlbumPhotos) // 获取相册照片列表

		// 需要认证的路由
		authGroup := albumGroup.Group("")
		authGroup.Use(authMiddleware.Handle())
		{
			authGroup.POST("/", h.CreateAlbum)                               // 创建相册
			authGroup.PUT("/:id", h.UpdateAlbum)                             // 更新相册
			authGroup.DELETE("/:id", h.DeleteAlbum)                          // 删除相册
			authGroup.POST("/:id/photos", h.AddPhotosToAlbum)                // 添加照片到相册
			authGroup.PUT("/:id/cover", h.SetAlbumCover)                     // 设置相册封面
			authGroup.DELETE("/:id/photos/:photoId", h.RemovePhotoFromAlbum) // 从相册删除照片
		}
	}
}

// ListAlbums 获取相册列表
// @Summary 获取相册列表
// @Description 获取相册列表，保持分页数据结构
// @Tags albums
// @Accept json
// @Produce json
// @Param page query int false "页码，默认1"
// @Param size query int false "每页数量，默认10"
// @Success 200 {object} Response{data=service.AlbumListResponse}
// @Router /api/v1/albums [get]
func (h *AlbumHandler) ListAlbums(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 调用服务层获取相册列表
	albums, err := h.AlbumService.ListAlbums(ctx, page, size)
	if err != nil {
		h.AlbumService.Log.Error("获取相册列表失败", "error", err, "page", page, "size", size)
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
		Data:    albums,
	})
}

// ListAlbumPhotos 获取相册照片列表
// @Summary 获取相册照片列表
// @Description 获取指定相册的照片列表，保持分页数据结构
// @Tags albums
// @Accept json
// @Produce json
// @Param id path int true "相册ID"
// @Param page query int false "页码，默认1"
// @Param size query int false "每页数量，默认10"
// @Success 200 {object} Response{data=service.AlbumPhotoListResponse}
// @Router /api/v1/albums/{id}/photos [get]
func (h *AlbumHandler) ListAlbumPhotos(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析相册ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.AlbumService.Log.Error("无效的相册ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的相册ID",
			Data:    nil,
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 调用服务层获取相册照片列表
	photos, err := h.AlbumService.ListAlbumPhotos(ctx, id, page, size)
	if err != nil {
		h.AlbumService.Log.Error("获取相册照片列表失败", "albumId", id, "error", err, "page", page, "size", size)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if photos == nil {
		h.AlbumService.Log.Info("相册不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "相册不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "查询成功",
		Data:    photos,
	})
}

// CreateAlbum 创建相册
// @Summary 创建相册
// @Description 创建一个新的相册
// @Tags albums
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param album body service.AlbumCreateRequest true "相册信息"
// @Success 200 {object} Response{data=service.Album}
// @Router /api/v1/albums [post]
func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析请求参数
	var req service.AlbumCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.AlbumService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 调用服务层创建相册
	album, err := h.AlbumService.CreateAlbum(ctx, &req)
	if err != nil {
		h.AlbumService.Log.Error("创建相册失败", "error", err, "request", req)
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
		Data:    album,
	})
}

// UpdateAlbum 更新相册
// @Summary 更新相册
// @Description 更新相册的信息
// @Tags albums
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "相册ID"
// @Param album body service.AlbumUpdateRequest true "相册信息"
// @Success 200 {object} Response{data=service.Album}
// @Router /api/v1/albums/{id} [put]
func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析相册ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.AlbumService.Log.Error("无效的相册ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的相册ID",
			Data:    nil,
		})
		return
	}

	// 解析请求参数
	var req service.AlbumUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.AlbumService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 调用服务层更新相册
	album, err := h.AlbumService.UpdateAlbum(ctx, id, &req)
	if err != nil {
		h.AlbumService.Log.Error("更新相册失败", "id", id, "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if album == nil {
		h.AlbumService.Log.Info("相册不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "相册不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "更新成功",
		Data:    album,
	})
}

// DeleteAlbum 删除相册
// @Summary 删除相册
// @Description 删除指定的相册
// @Tags albums
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "相册ID"
// @Success 200 {object} Response{data=nil}
// @Router /api/v1/albums/{id} [delete]
func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析相册ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.AlbumService.Log.Error("无效的相册ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的相册ID",
			Data:    nil,
		})
		return
	}

	// 调用服务层删除相册
	ok, err := h.AlbumService.DeleteAlbum(ctx, id)
	if err != nil {
		h.AlbumService.Log.Error("删除相册失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if !ok {
		h.AlbumService.Log.Info("相册不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "相册不存在",
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

// AddPhotosToAlbum 更新相册照片
// @Summary 更新相册照片
// @Description 更新相册的照片列表，只会保留指定的照片ID
// @Tags albums
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "相册ID"
// @Param photos body service.AlbumAddPhotosRequest true "照片ID列表"
// @Success 200 {object} Response{data=[]service.AlbumPhoto}
// @Router /api/v1/albums/{id}/photos [post]
func (h *AlbumHandler) AddPhotosToAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析相册ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.AlbumService.Log.Error("无效的相册ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的相册ID",
			Data:    nil,
		})
		return
	}

	// 解析请求参数
	var req service.AlbumAddPhotosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.AlbumService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 调用服务层添加照片到相册
	photos, err := h.AlbumService.AddPhotosToAlbum(ctx, id, &req)
	if err != nil {
		h.AlbumService.Log.Error("添加照片到相册失败", "albumId", id, "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if photos == nil {
		h.AlbumService.Log.Info("相册不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "相册不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "添加成功",
		Data:    photos,
	})
}

// SetAlbumCover 设置相册封面
// @Summary 设置相册封面
// @Description 设置指定相册的封面照片
// @Tags albums
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "相册ID"
// @Param cover body service.AlbumSetCoverRequest true "封面照片ID"
// @Success 200 {object} Response{data=service.Album}
// @Router /api/v1/albums/{id}/cover [put]
func (h *AlbumHandler) SetAlbumCover(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析相册ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.AlbumService.Log.Error("无效的相册ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的相册ID",
			Data:    nil,
		})
		return
	}

	// 解析请求参数
	var req service.AlbumSetCoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.AlbumService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "请求参数错误: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 调用服务层设置相册封面
	album, err := h.AlbumService.SetAlbumCover(ctx, id, &req)
	if err != nil {
		h.AlbumService.Log.Error("设置相册封面失败", "albumId", id, "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if album == nil {
		h.AlbumService.Log.Info("相册不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "相册不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "设置成功",
		Data:    album,
	})
}

// RemovePhotoFromAlbum 从相册删除照片
// @Summary 从相册删除照片
// @Description 从指定的相册删除照片
// @Tags albums
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "相册ID"
// @Param photoId path int true "照片ID"
// @Success 200 {object} Response{data=nil}
// @Router /api/v1/albums/{id}/photos/{photoId} [delete]
func (h *AlbumHandler) RemovePhotoFromAlbum(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析相册ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.AlbumService.Log.Error("无效的相册ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的相册ID",
			Data:    nil,
		})
		return
	}

	// 解析照片ID
	photoIdStr := c.Param("photoId")
	photoId, err := strconv.ParseUint(photoIdStr, 10, 64)
	if err != nil {
		h.AlbumService.Log.Error("无效的照片ID", "photoId", photoIdStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的照片ID",
			Data:    nil,
		})
		return
	}

	// 调用服务层从相册删除照片
	ok, err := h.AlbumService.RemovePhotoFromAlbum(ctx, id, photoId)
	if err != nil {
		h.AlbumService.Log.Error("从相册删除照片失败", "albumId", id, "photoId", photoId, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if !ok {
		h.AlbumService.Log.Info("相册不存在", "id", id)
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "相册不存在",
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
