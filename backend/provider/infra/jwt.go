package infra

import (
	"github.com/bookandmusic/love-girl/internal/auth"
	"github.com/bookandmusic/love-girl/internal/config"
)

func ProvideJWT(cfg *config.AppConfig) auth.JWT {
	return auth.NewHS256JWT(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.Expire)
}
