package infra

import (
	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/log"
)

func ProvideLogger(cfg *config.AppConfig) *log.Logger {
	return log.NewLogger(cfg.Log)
}

func ProvideGormLogger(
	cfg *config.AppConfig,
	logger *log.Logger,
) *log.GormLogger {
	return log.NewGormLogger(logger, cfg.Log)
}
