// package log

package log

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"

	"github.com/bookandmusic/love-girl/internal/config"
)

// Logger 封装 slog.Logger，便于扩展和实现 io.Writer
type Logger struct {
	*slog.Logger
}

// Init 初始化 slog 日志（支持 JSON 和 Text 格式）
func getLoggerLevel(level string) slog.Level {
	// 设置日志级别
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return logLevel
}

func NewLogger(cfg config.LogConfig) *Logger {
	level := getLoggerLevel(cfg.Level)
	opts := &slog.HandlerOptions{Level: level, AddSource: true}

	var handler slog.Handler

	// 根据配置选择格式
	if cfg.Format == "json" {
		// JSON 格式 - 用于日志收集
		var writer io.Writer = os.Stdout
		if cfg.Output != "" && cfg.Output != "stdout" {
			file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err == nil {
				writer = file
			}
		}
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		// 彩色文本格式 - 默认行为
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			AddSource:  true,
			TimeFormat: "2006-01-02 15:04:05",
		})
	}

	logger := slog.New(handler)

	// 设置为全局 logger（推荐）
	slog.SetDefault(logger)

	return &Logger{logger}
}

// Write 实现 io.Writer 接口，用于 Gin 等框架的日志输出
func (l *Logger) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg != "" {
		l.Info(msg)
	}
	return len(p), nil
}
