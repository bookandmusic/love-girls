package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 错误恢复中间件（增强版）
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())
				requestID := c.GetString(RequestIDKey)

				slog.Error("panic recovered",
					"request_id", requestID,
					"error", err,
					"stack", stack,
				)

				// 生产环境不暴露详细错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    1,
					"message": "系统内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
