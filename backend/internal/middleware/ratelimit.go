package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IPRateLimiter IP 限流器
type IPRateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // 时间窗口内允许的最大请求数
	window   time.Duration // 时间窗口
}

type visitor struct {
	windowStart time.Time // 窗口开始时间
	count       int       // 当前窗口内的请求计数
}

// NewIPRateLimiter 创建新的 IP 限流器
func NewIPRateLimiter(rate int, window time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// 启动后台清理协程，定期清理过期的访问者记录
	go limiter.cleanupVisitors()

	return limiter
}

// Allow 检查是否允许请求
func (l *IPRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	v, exists := l.visitors[ip]

	// 不存在或窗口已过期，创建新窗口
	if !exists || now.Sub(v.windowStart) > l.window {
		l.visitors[ip] = &visitor{
			windowStart: now,
			count:       1,
		}
		return true
	}

	// 在当前窗口内，检查是否超过限制
	if v.count < l.rate {
		v.count++
		return true
	}

	// 超过限制
	return false
}

// cleanupVisitors 定期清理过期的访问者记录
func (l *IPRateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for ip, v := range l.visitors {
			if now.Sub(v.windowStart) > l.window {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

// RateLimit 返回 IP 限流中间件
// rate: 时间窗口内允许的最大请求数
// window: 时间窗口
func RateLimit(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    1,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByKey 返回基于自定义 key 的限流中间件
// keyFunc: 自定义 key 生成函数
// rate: 时间窗口内允许的最大请求数
// window: 时间窗口
func RateLimitByKey(keyFunc func(c *gin.Context) string, rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate, window)

	return func(c *gin.Context) {
		key := keyFunc(c)

		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    1,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
