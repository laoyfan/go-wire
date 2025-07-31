package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-wire/config"
	"go-wire/logger"
	"go-wire/redis"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Engine *gin.Engine
	Config *config.Config
	Logger logger.Logger
	Redis  *redis.RedisClients
}

var ProviderSet = wire.NewSet(NewApp)

func NewApp(engine *gin.Engine, cfg *config.Config, log logger.Logger, redisClients *redis.RedisClients) *App {
	return &App{
		Engine: engine,
		Config: cfg,
		Logger: log,
		Redis:  redisClients,
	}
}

func (app *App) Run() error {
	// 1. Gin 基础设置
	gin.DisableConsoleColor()
	gin.SetMode(app.Config.App.Mode)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", app.Config.App.Port),
		Handler:        app.Engine,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ctx, cancel := app.createContextWithTraceID()
	// 在服务关闭时断开 Redis 连接
	defer func() {
		if err := app.Redis.Close(ctx); err != nil {
			app.Logger.Error(ctx, "关闭 Redis 出错", logger.ErrorField(err))
		}
	}()
	defer cancel()
	go func() {
		app.startServer(ctx, server)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.Logger.Info(ctx, "关闭服务...")
	// 关闭服务
	if err := server.Shutdown(ctx); err != nil {
		app.Logger.Fatal(ctx, "服务关闭原因:", logger.ErrorField(err))
	}
	app.Logger.Info(ctx, "服务退出")
	return nil
}

// 创建包含 Trace ID 的上下文
func (app *App) createContextWithTraceID() (context.Context, context.CancelFunc) {
	baseCtx := context.Background()
	traceID := fmt.Sprintf("main:date(%s)", time.Now().Format("2006-01-02 15:04:05"))
	ctx := context.WithValue(baseCtx, "TraceID", traceID)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	return ctx, cancel
}

// 启动 HTTP 服务器
func (app *App) startServer(ctx context.Context, server *http.Server) {
	app.Logger.Info(ctx, fmt.Sprintf("服务开启:%d", app.Config.App.Port))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.Logger.Fatal(ctx, "listen: %s\n", logger.ErrorField(err))
	}
}
