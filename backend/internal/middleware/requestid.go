package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDKey 是 Context 中存储 RequestID 的键
	RequestIDKey = "request_id"
	// RequestIDHeader 是 HTTP 头部名称
	RequestIDHeader = "X-Request-ID"
)

// RequestID 返回一个中间件，为每个请求生成或传递 RequestID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先使用上游传递的 RequestID
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 存储到 Context 和 Response Header
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}
