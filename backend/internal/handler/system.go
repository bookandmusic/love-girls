package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/internal/service"
)

type SystemHandler struct {
	SystemService *service.SystemService
}

func NewSystemHandler(systemService *service.SystemService) *SystemHandler {
	return &SystemHandler{
		SystemService: systemService,
	}
}

// RegisterRoutes 注册系统相关的路由
func (h *SystemHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	systemGroup := apiGroup.Group("/system")
	{
		// 系统初始化相关路由
		systemGroup.POST("/init", h.InitSystem)
		systemGroup.GET("/init", h.CheckSystemInitStatus)

		// 系统信息相关路由
		systemGroup.GET("/info", h.GetSystemInfo)

		// 站点设置相关路由 (需要认证)
		settingsGroup := systemGroup.Group("/settings")
		settingsGroup.Use(authMiddleware.Handle())
		{
			settingsGroup.POST("/site", h.SaveSettings)
			settingsGroup.GET("/site", h.GetSettings)
		}

		// 仪表盘统计数据路由 (需要认证)
		dashboardGroup := systemGroup.Group("/dashboard")
		dashboardGroup.Use(authMiddleware.Handle())
		{
			dashboardGroup.GET("/stats", h.GetDashboardStats)
		}
	}
}

// InitSystem godoc
// @Summary 初始化系统
// @Description 初始化系统，包括创建站点信息、用户信息和设置访问密码
// @Tags system
// @Accept json
// @Produce json
// @Param initSystem body service.InitSystemRequest true "系统初始化参数"
// @Success 200 {object} Response{data=map[string]bool}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /system/init [post]
func (h *SystemHandler) InitSystem(c *gin.Context) {
	ctx := c.Request.Context()

	var req service.InitSystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SystemService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "参数校验失败",
			Data:    nil,
		})
		return
	}

	if err := h.SystemService.InitSystem(ctx, &req); err != nil {
		h.SystemService.Log.Error("系统初始化失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统初始化失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "系统初始化成功",
		Data: map[string]bool{
			"initialized": true,
		},
	})
}

// CheckSystemInitStatus godoc
// @Summary 检查系统初始化状态
// @Description 检查系统是否已完成初始化
// @Tags system
// @Produce json
// @Success 200 {object} Response{data=map[string]bool}
// @Failure 500 {object} Response
// @Router /system/init [get]
func (h *SystemHandler) CheckSystemInitStatus(c *gin.Context) {
	ctx := c.Request.Context()

	initialized, err := h.SystemService.IsSystemInitialized(ctx)
	if err != nil {
		h.SystemService.Log.Error("检查初始化状态失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "检查初始化状态失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "查询成功",
		Data: map[string]bool{
			"initialized": initialized,
		},
	})
}

// GetSystemInfo godoc
// @Summary 获取系统信息
// @Description 获取系统配置的基本信息
// @Tags system
// @Produce json
// @Success 200 {object} Response{data=service.FrontendSystemInfo}
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /system/info [get]
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	ctx := c.Request.Context()

	info, err := h.SystemService.GetSystemInfo(ctx)
	if err != nil {
		h.SystemService.Log.Error("获取系统信息失败", "error", err)
		if err.Error() == "系统未初始化" {
			c.JSON(http.StatusNotFound, Response{
				Code:    1,
				Message: "系统未初始化或获取信息失败",
				Data:    nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    1,
				Message: "系统内部错误",
				Data:    nil,
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "查询成功",
		Data:    info,
	})
}

// GetSettings godoc
// @Summary 获取站点设置
// @Description 获取当前站点的设置信息
// @Security OAuth2Password
// @Tags system
// @Produce json
// @Success 200 {object} Response{data=map[string]string}
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /system/settings/site [get]
func (h *SystemHandler) GetSettings(c *gin.Context) {
	ctx := c.Request.Context()

	settings, err := h.SystemService.GetSettings(ctx)
	if err != nil {
		h.SystemService.Log.Error("获取站点设置失败", "error", err)
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
		Data:    settings,
	})
}

// SaveSettings godoc
// @Summary 保存站点设置
// @Description 更新站点的设置信息
// @Security OAuth2Password
// @Tags system
// @Accept json
// @Produce json
// @Param settings body map[string]string true "站点设置"
// @Success 200 {object} Response{data=map[string]string}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /system/settings/site [post]
func (h *SystemHandler) SaveSettings(c *gin.Context) {
	ctx := c.Request.Context()

	var settings map[string]string
	if err := c.ShouldBindJSON(&settings); err != nil {
		h.SystemService.Log.Error("参数校验失败", "error", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    1,
			Message: "参数校验失败",
			Data:    nil,
		})
		return
	}

	if err := h.SystemService.SaveSettings(ctx, settings); err != nil {
		h.SystemService.Log.Error("保存站点设置失败", "error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    1,
			Message: "系统内部错误",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "保存成功",
		Data:    settings,
	})
}

// GetDashboardStats godoc
// @Summary 获取仪表盘统计数据
// @Description 获取仪表盘的统计数据，用于展示系统的整体运营情况和数据概览
// @Security OAuth2Password
// @Tags system
// @Produce json
// @Success 200 {object} Response{data=model.DashboardStats}
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /system/dashboard/stats [get]
func (h *SystemHandler) GetDashboardStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.SystemService.GetDashboardStats(ctx)
	if err != nil {
		h.SystemService.Log.Error("获取仪表盘统计数据失败", "error", err)
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
		Data:    stats,
	})
}
