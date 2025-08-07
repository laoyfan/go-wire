package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"go-wire/config"
	"go-wire/logger"
	"go-wire/redis"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type App struct {
	engine *gin.Engine
	cfg    *config.Config
	redis  *redis.Redis
	log    logger.Logger
}

var ProviderSet = wire.NewSet(NewApp)

func NewApp(engine *gin.Engine, cfg *config.Config, redis *redis.Redis, log logger.Logger) *App {
	return &App{
		engine: engine,
		cfg:    cfg,
		redis:  redis,
		log:    log,
	}
}

func (app *App) Run() error {
	// 1. Gin 基础设置
	gin.DisableConsoleColor()
	gin.SetMode(app.cfg.App.Mode)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", app.cfg.App.Port),
		Handler:        app.engine,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ctx := context.WithValue(context.Background(), "TraceID", fmt.Sprintf("main:date:%s", time.Now().Format("2006-01-02 15:04:05")))
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 在服务关闭时断开 Redis 连接
	defer func() {
		if err := app.redis.Close(shutdownCtx); err != nil {
			app.log.Error(ctx, "关闭 Redis 出错", logger.Error(err))
		}
		_ = app.log.Sync()
	}()
	go func() {
		app.log.Info(ctx, fmt.Sprintf("服务开启:%d", app.cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.log.Fatal(ctx, "listen: %s\n", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	app.log.Info(ctx, fmt.Sprintf("收到退出信号: %v", sig))
	// 关闭服务
	if err := server.Shutdown(shutdownCtx); err != nil {
		app.log.Fatal(ctx, "服务关闭原因:", logger.Error(err))
	}
	app.log.Info(ctx, "服务退出")
	return nil
}
