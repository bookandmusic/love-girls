package infra

import "github.com/bookandmusic/love-girl/internal/config"

func ProvideConfig() (*config.AppConfig, error) {
	return config.LoadConfig()
}
