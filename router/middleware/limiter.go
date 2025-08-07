package middleware

import (
	"go-wire/config"
	"go-wire/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type LimiterMiddleware struct {
	log     logger.Logger
	limiter *rate.Limiter
}

func NewLimiterMiddleware(log logger.Logger, limiter *rate.Limiter) *LimiterMiddleware {
	return &LimiterMiddleware{
		log:     log,
		limiter: limiter,
	}
}

func (m *LimiterMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 检查请求是否被限流
		if !m.limiter.Allow() {
			m.log.Warn(ctx, "请求被限流",
				logger.StringAny("url", ctx.Request.URL.Path),
				logger.StringAny("client_ip", ctx.ClientIP()),
			)
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": http.StatusTooManyRequests,
				"msg":  "服务繁忙，请稍后再试...",
			})
			return
		}
		ctx.Next()
	}
}

func NewLimiter(cfg *config.Config) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(cfg.App.Limit), cfg.App.Burst)
}
