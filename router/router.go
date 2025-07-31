package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewRouter)

func NewRouter() *gin.Engine {
	// 创建 gin.Engine 并设置模式
	engine := gin.New()

	// 可以根据 cfg 做条件配置，如是否开启 swagger、debug 日志中间件等

	// 注册中间件 + 路由
	registerMiddleware(engine)
	registerRoutes(engine)

	return engine
}

func registerMiddleware(r *gin.Engine) {
	// 例如 Logger、Recovery、中间件顺序不要乱
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
}

func registerRoutes(r *gin.Engine) {
	// 实际业务路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	RegisterApi(r)
}
