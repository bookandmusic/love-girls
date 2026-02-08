package infra

import (
	"github.com/bookandmusic/love-girl/internal/auth"
	"github.com/bookandmusic/love-girl/internal/middleware"
)

func ProvideAuthMiddleware(jwt auth.JWT) *middleware.AuthMiddleware {
	return middleware.NewAuthMiddleware(jwt)
}
