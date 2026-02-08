package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/log"
)

// ErrRestart 表示服务因配置变更需要重启
var ErrRestart = errors.New("server restart requested")

type App struct {
	Logger     *log.Logger
	Engine     *gin.Engine
	App        config.AppConfigApp
	Server     config.ServerConfig
	httpServer *http.Server
}

func NewApp(
	logger *log.Logger,
	engine *gin.Engine,
	cfg config.AppConfig,
) *App {
	app := &App{
		Logger: logger,
		Engine: engine,
		Server: cfg.Server,
		App:    cfg.App,
	}

	app.httpServer = &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return app
}

func (a *App) Run(
	version, commit, buildTime string,
	restartCh chan struct{},
) error {
	a.Logger.Info(fmt.Sprintf("Server %s, Starting %s, mode: %s, version: %s commit: %s, buildTime: %s",
		a.App.Name, a.Server.Addr, a.Server.Mode, version, commit, buildTime))

	// 启动配置文件监听
	v := config.GetViperInstance()
	if v != nil {
		configPath := "./data/configs/config.yaml"
		config.WatchConfig(v, configPath, restartCh)
		a.Logger.Info("配置文件监听已启动")
	} else {
		a.Logger.Warn("未找到 viper 实例，配置文件监听未启动")
	}

	// 监听重启信号，收到后优雅关闭 HTTP server
	restart := false
	go func() {
		<-restartCh
		restart = true
		a.Logger.Info("收到重启信号，正在优雅关闭服务器...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := a.httpServer.Shutdown(ctx); err != nil {
			a.Logger.Error("服务器关闭错误", "error", err)
		}
	}()

	a.Logger.Info("服务已就绪，支持配置热加载")

	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.Logger.Error("服务器运行错误", "error", err)
		return err
	}

	if restart {
		a.Logger.Info("服务器已停止，准备重启...")
		return ErrRestart
	}

	a.Logger.Info("服务器已停止")
	return nil
}

func (a *App) GracefulShutdown() {
	a.Logger.Info("正在优雅关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.Logger.Error("服务器关闭错误", "error", err)
	}
	a.Logger.Info("服务器已优雅关闭")
}
