package log

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"

	"github.com/bookandmusic/love-girl/internal/config"
)

// GormLogger 实现 gorm/logger.Interface
type GormLogger struct {
	Logger *Logger
	Level  logger.LogLevel
}

// NewGormLogger 创建 GORM Logger
func NewGormLogger(logger *Logger, cfg config.LogConfig) *GormLogger {
	level := getGormLoggerLevel(cfg.Level)
	return &GormLogger{
		Logger: logger,
		Level:  level,
	}
}

func getGormLoggerLevel(level string) logger.LogLevel {
	// 设置日志级别
	var logLevel logger.LogLevel
	switch level {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	default:
		logLevel = logger.Info
	}

	return logLevel
}

// LogMode 设置日志等级
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.Level = level
	return &newLogger
}

// Info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.Level >= logger.Info {
		l.Logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.Level >= logger.Warn {
		l.Logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error
func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.Level >= logger.Error {
		l.Logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	msg := fmt.Sprintf("[%s] [rows:%d] %s", elapsed, rows, sql)
	if err != nil {
		l.Logger.Error(msg, "error", err)
		return
	}

	switch {
	case l.Level >= logger.Info && elapsed > 500*time.Millisecond:
		l.Logger.Warn(msg)
	case l.Level >= logger.Info:
		l.Logger.Info(msg)
	}
}
