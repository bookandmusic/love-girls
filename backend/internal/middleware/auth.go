package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/love-girl/internal/auth"
)

func NewAuthMiddleware(jwt auth.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		JWT: jwt,
	}
}

type AuthMiddleware struct {
	JWT auth.JWT
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractBearerToken(c)
		if token == "" {
			m.unauthorized(c)
			return
		}

		claims, err := m.JWT.Parse(token)
		if err != nil {
			m.unauthorized(c)
			return
		}

		// 注入 Claims，不查 DB
		auth.SetAuthClaims(c, claims)

		c.Next()
	}
}

func (m *AuthMiddleware) unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
	})
}

func (m *AuthMiddleware) extractBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
