package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout 超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if timeout <= 0 {
			c.Next()
			return
		}

		// 创建带超时的 Context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替换请求 Context
		c.Request = c.Request.WithContext(ctx)

		// 使用 channel 监听完成
		finished := make(chan struct{})
		go func() {
			defer close(finished)
			c.Next()
		}()

		select {
		case <-finished:
			// 请求正常完成
		case <-ctx.Done():
			// 请求超时
			c.JSON(http.StatusRequestTimeout, gin.H{
				"code":    1,
				"message": "请求超时",
			})
			c.Abort()
		}
	}
}
