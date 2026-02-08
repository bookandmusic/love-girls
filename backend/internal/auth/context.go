package auth

import (
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserID uint64
	Role   string
}

type contextKey string

const authContextKey contextKey = "authClaims"

// SetAuthClaims 写入认证信息（仅中间件使用）
func SetAuthClaims(c *gin.Context, claims *Claims) {
	c.Set(string(authContextKey), claims)
}

// GetAuthClaims 获取认证信息
func GetAuthClaims(c *gin.Context) (*Claims, bool) {
	v, ok := c.Get(string(authContextKey))
	if !ok {
		return nil, false
	}
	claims, ok := v.(*Claims)
	return claims, ok
}

// MustGetAuthClaims 强制获取（已通过 Auth Middleware 的接口可用）
func MustGetAuthClaims(c *gin.Context) *Claims {
	claims, ok := GetAuthClaims(c)
	if !ok {
		panic("auth claims not found in context")
	}
	return claims
}
