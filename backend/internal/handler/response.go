package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError 应用层错误
type AppError struct {
	Code    int    // HTTP 状态码
	Message string // 用户可见消息
	Detail  string // 详细错误信息（可选）
	Cause   error  // 原始错误（可选）
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// 预定义错误
var (
	ErrBadRequest   = &AppError{Code: http.StatusBadRequest, Message: "请求参数错误"}
	ErrUnauthorized = &AppError{Code: http.StatusUnauthorized, Message: "未授权访问"}
	ErrForbidden    = &AppError{Code: http.StatusForbidden, Message: "禁止访问"}
	ErrNotFound     = &AppError{Code: http.StatusNotFound, Message: "资源不存在"}
	ErrInternal     = &AppError{Code: http.StatusInternalServerError, Message: "系统内部错误"}
)

// NewAppError 创建应用错误
func NewAppError(code int, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, Response{
			Code:    1,
			Message: appErr.Message,
		})
		return
	}

	// 未知错误，返回 500
	c.JSON(http.StatusInternalServerError, Response{
		Code:    1,
		Message: "系统内部错误",
	})
}

// FailWithDetail 失败响应（带详情）
func FailWithDetail(c *gin.Context, err error, detail string) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		appErr.Detail = detail
		c.JSON(appErr.Code, Response{
			Code:    1,
			Message: appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code:    1,
		Message: "系统内部错误",
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    1,
		Message: message,
	})
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    1,
		Message: "未授权访问",
	})
}

// Forbidden 403 错误
func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code:    1,
		Message: "禁止访问",
	})
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    1,
		Message: message,
	})
}

// InternalError 500 错误
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "系统内部错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    1,
		Message: message,
	})
}
