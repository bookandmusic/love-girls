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

type MomentHandler struct {
	MomentService *service.MomentService
}

func NewMomentHandler(momentService *service.MomentService) *MomentHandler {
	return &MomentHandler{
		MomentService: momentService,
	}
}

// RegisterRoutes 注册动态相关的路由
func (h *MomentHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	apiGroup.GET("/moments", h.ListMoments) // 获取动态列表（公开的）
	// 需要认证的路由
	authGroup := apiGroup.Group("")
	authGroup.Use(authMiddleware.Handle())
	{
		authGroup.POST("/moments", h.CreateMoment)           // 创建动态
		authGroup.PUT("/moments/:id", h.UpdateMoment)        // 更新动态
		authGroup.DELETE("/moments/:id", h.DeleteMoment)     // 删除动态
		authGroup.PUT("/moments/:id/public", h.UpdatePublic) // 更新动态公开状态
		authGroup.POST("/moments/:id/like", h.LikeMoment)    // 点赞动态
	}

}

// CreateMoment 创建动态
// @Summary 创建动态
// @Description 创建一个新的动态
// @Tags moments
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param moment body service.MomentCreateRequest true "动态信息"
// @Success 200 {object} Response{data=service.FrontendMoment}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /moments [post]
func (h *MomentHandler) CreateMoment(c *gin.Context) {
	ctx := c.Request.Context()

	var req service.MomentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.MomentService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "参数校验失败",
			Data:    nil,
		})
		return
	}

	moment, err := h.MomentService.CreateMoment(ctx, &req)
	if err != nil {
		h.MomentService.Log.Error("创建动态失败", "error", err)
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
		Data:    moment,
	})
}

// ListMoments 获取动态列表
// @Summary 获取动态列表
// @Description 分页获取动态列表
// @Tags moments
// @Produce json
// @Param page query int true "页码"
// @Param size query int true "每页数量"
// @Success 200 {object} Response{data=service.MomentListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /moments [get]
func (h *MomentHandler) ListMoments(c *gin.Context) {
	ctx := c.Request.Context()

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.MomentService.Log.Error("无效的页码参数", "page", pageStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的页码参数",
			Data:    nil,
		})
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		h.MomentService.Log.Error("无效的每页数量参数", "size", sizeStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的每页数量参数",
			Data:    nil,
		})
		return
	}

	// 检查用户是否已登录
	claims, isLoggedIn := auth.GetAuthClaims(c)

	var userID uint64
	if isLoggedIn {
		userID = claims.UserID
	}

	// 根据登录状态决定是否获取所有动态或仅公开动态
	listResp, err := h.MomentService.ListMomentsByAuthStatus(ctx, page, size, isLoggedIn, userID)
	if err != nil {
		h.MomentService.Log.Error("获取动态列表失败", "error", err)
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
		Data:    listResp,
	})
}

// UpdateMoment 更新动态
// @Summary 更新动态
// @Description 更新动态信息
// @Tags moments
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "动态ID"
// @Param moment body service.MomentUpdateRequest true "动态信息"
// @Success 200 {object} Response{data=service.FrontendMoment}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /moments/{id} [put]
func (h *MomentHandler) UpdateMoment(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.MomentService.Log.Error("无效的动态ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的动态ID",
			Data:    nil,
		})
		return
	}

	var req service.MomentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.MomentService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "参数校验失败",
			Data:    nil,
		})
		return
	}

	moment, err := h.MomentService.UpdateMoment(ctx, id, &req)
	if err != nil {
		h.MomentService.Log.Error("更新动态失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if moment == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "动态不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "更新成功",
		Data:    moment,
	})
}

// DeleteMoment 删除动态
// @Summary 删除动态
// @Description 删除指定ID的动态
// @Tags moments
// @Produce json
// @Security OAuth2Password
// @Param id path string true "动态ID"
// @Success 200 {object} Response{data=interface{}}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /moments/{id} [delete]
func (h *MomentHandler) DeleteMoment(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.MomentService.Log.Error("无效的动态ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的动态ID",
			Data:    nil,
		})
		return
	}

	deleted, err := h.MomentService.DeleteMoment(ctx, id)
	if err != nil {
		h.MomentService.Log.Error("删除动态失败", "id", id, "error", err)
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
			Message: "动态不存在",
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

// UpdatePublic 更新动态公开状态
// @Summary 更新动态公开状态
// @Description 更新动态的公开状态
// @Tags moments
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "动态ID"
// @Param status body service.MomentPublicRequest true "公开状态"
// @Success 200 {object} Response{data=service.FrontendMoment}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /moments/{id}/public [put]
func (h *MomentHandler) UpdatePublic(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.MomentService.Log.Error("无效的动态ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的动态ID",
			Data:    nil,
		})
		return
	}

	var req service.MomentPublicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.MomentService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "参数校验失败",
			Data:    nil,
		})
		return
	}

	moment, err := h.MomentService.UpdatePublicStatus(ctx, id, req.IsPublic)
	if err != nil {
		h.MomentService.Log.Error("更新动态公开状态失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if moment == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "动态不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "更新成功",
		Data:    moment,
	})
}

// LikeMoment 点赞动态
// @Summary 点赞动态
// @Description 为指定ID的动态点赞
// @Tags moments
// @Produce json
// @Security OAuth2Password
// @Param id path string true "动态ID"
// @Success 200 {object} Response{data=service.MomentLikeResponse}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /moments/{id}/like [post]
func (h *MomentHandler) LikeMoment(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.MomentService.Log.Error("无效的动态ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "无效的动态ID",
			Data:    nil,
		})
		return
	}

	liked, err := h.MomentService.LikeMoment(ctx, id)
	if err != nil {
		h.MomentService.Log.Error("点赞动态失败", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	if liked == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    1,
			Message: "动态不存在",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "点赞成功",
		Data: service.MomentLikeResponse{
			Likes: liked.Likes,
		},
	})
}
