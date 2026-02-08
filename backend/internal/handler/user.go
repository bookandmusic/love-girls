package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/love-girl/internal/auth"
	middle "github.com/bookandmusic/love-girl/internal/middleware"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/internal/service"
)

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

type UserHandler struct {
	UserService *service.UserService
	JWT         auth.JWT
}

// RegisterRoutes 注册用户相关的路由
func (h *UserHandler) RegisterRoutes(apiGroup *gin.RouterGroup, server *server.GinEngine, authMiddleware *middle.AuthMiddleware) {
	// 登录接口
	apiGroup.POST("/user/token", h.Login)

	// 需要认证的路由
	authGroup := apiGroup.Group("")
	authGroup.Use(authMiddleware.Handle())
	{
		// 用户信息接口
		authGroup.GET("/user", h.GetUserInfo)
		// 用户管理接口
		authGroup.GET("/users", h.GetUsers)
		authGroup.GET("/users/:id/avatars", h.GetUserAvatarHistory)
		authGroup.PUT("/users/:id", h.UpdateUser)
	}
}

type UserLoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// @Summary Generate User Token
// @Description User login with username and password to get token
// @Tags user
// @Accept json
// @Produce json
// @Param user body UserLoginRequest true "User login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /user/token [post]
func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req UserLoginRequest
	if err := c.ShouldBind(&req); err != nil {
		h.UserService.Log.Error("参数校验失败", "error", err)
		c.JSON(
			http.StatusBadRequest,
			Response{
				Code:    1,
				Message: "参数格式错误或字段缺失",
				Data:    nil,
			},
		)
		return
	}

	_, token, err := h.UserService.GenerateToken(ctx, req.Username, req.Password)
	if err != nil {
		h.UserService.Log.Error("用户登录失败", "error", err, "username", req.Username)
		c.JSON(
			http.StatusUnauthorized,
			Response{
				Code:    1,
				Message: "用户名或密码错误",
				Data:    nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"access_token": token,
			"token_type":   "bearer",
			"expires_in":   3600,
		},
	)
}

// @Summary Get user info
// @Description Get user info with token
// @Tags user
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /user [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	claims := auth.MustGetAuthClaims(c)
	user, err := h.UserService.GetUserInfo(ctx, claims.UserID)
	if err != nil {
		h.UserService.Log.Error("用户信息查询失败", "error", err, "userID", claims.UserID)
		c.JSON(
			http.StatusInternalServerError,
			Response{
				Code:    1,
				Message: "系统内部错误",
				Data:    nil,
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		Response{
			Code:    0,
			Message: "查询成功",
			Data: map[string]interface{}{
				"userName":  user.Name,
				"userId":    user.ID,
				"userEmail": getEmailValue(user.Email),
			},
		},
	)
}

// getEmailValue 获取Email的实际值，处理指针类型
func getEmailValue(email *string) string {
	if email != nil {
		return *email
	}
	return ""
}

// @Summary Get users list
// @Description Get all users list
// @Tags user
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.UserService.GetUsers(ctx)
	if err != nil {
		h.UserService.Log.Error("获取用户列表失败", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			Response{
				Code:    1,
				Message: "系统内部错误",
				Data:    nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		Response{
			Code:    0,
			Message: "查询成功",
			Data:    users,
		},
	)
}

// @Summary Update user info
// @Description Update user info
// @Tags user
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path int true "User ID"
// @Param user body map[string]interface{} true "User info"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	// 解析路径参数
	var userID uint64
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &userID); err != nil {
		c.JSON(
			http.StatusBadRequest,
			Response{
				Code:    1,
				Message: "用户ID格式错误",
				Data:    nil,
			},
		)
		return
	}

	// 解析请求体
	type UpdateUserRequest struct {
		Name        string  `json:"name" binding:"required"`
		Email       string  `json:"email" binding:"omitempty,email"`
		Avatar      string  `json:"avatar"`
		AvatarID    *uint64 `json:"avatarId"`
		Role        string  `json:"role"`
		NewPassword string  `json:"newPassword"`
	}

	var req UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			Response{
				Code:    1,
				Message: "参数格式错误或字段缺失",
				Data:    nil,
			},
		)
		return
	}

	// 更新用户信息
	user, err := h.UserService.UpdateUser(ctx, userID, req.Name, req.Email, req.AvatarID, req.NewPassword)
	if err != nil {
		h.UserService.Log.Error("更新用户信息失败", "error", err, "userID", userID)
		c.JSON(
			http.StatusInternalServerError,
			Response{
				Code:    1,
				Message: "系统内部错误",
				Data:    nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		Response{
			Code:    0,
			Message: "更新成功",
			Data:    user,
		},
	)
}

// @Summary Get user avatar history
// @Description Get user avatar history with pagination
// @Tags user
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /users/{id}/avatars [get]
func (h *UserHandler) GetUserAvatarHistory(c *gin.Context) {
	ctx := c.Request.Context()

	var userID uint64
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &userID); err != nil {
		c.JSON(
			http.StatusBadRequest,
			Response{
				Code:    1,
				Message: "用户ID格式错误",
				Data:    nil,
			},
		)
		return
	}

	page, size := 1, 10
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			size = v
		}
	}

	result, err := h.UserService.GetAvatarHistory(ctx, userID, page, size)
	if err != nil {
		h.UserService.Log.Error("获取头像历史失败", "error", err, "userID", userID)
		c.JSON(
			http.StatusInternalServerError,
			Response{
				Code:    1,
				Message: "系统内部错误",
				Data:    nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		Response{
			Code:    0,
			Message: "查询成功",
			Data:    result,
		},
	)
}
