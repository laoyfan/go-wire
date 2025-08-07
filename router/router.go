package router

import (
	"go-wire/controller"
	"go-wire/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewRouter)

func NewRouter(
	auth *middleware.AuthMiddleware,
	cors *middleware.CorsMiddleware,
	trace *middleware.TraceMiddleware,
	limiter *middleware.LimiterMiddleware,
	logger *middleware.LoggerMiddleware,
	error *middleware.ErrorMiddleware,
	apiController *controller.ApiController,
) *gin.Engine {
	r := gin.New()
	// 注册所有中间件
	r.Use(auth.Handler())
	r.Use(cors.Handler())
	r.Use(trace.Handler())
	r.Use(limiter.Handler())
	r.Use(logger.Handler())
	r.Use(error.Handler())
	r.Use(trace.Handler())
	apiGroup := r.Group("/")
	apiController.RegisterRoutes(apiGroup)
	return r
}
