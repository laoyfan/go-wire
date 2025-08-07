package middleware

import (
	"go-wire/config"
	"go-wire/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct {
	cfg *config.Config
	log logger.Logger
}

func NewCorsMiddleware(cfg *config.Config, log logger.Logger) *CorsMiddleware {
	return &CorsMiddleware{cfg: cfg, log: log}
}

func (m *CorsMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")

		// 设置跨域响应头
		setHeaders(ctx, origin)

		// OPTIONS 方法直接返回
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 调试模式下放行所有请求
		if m.cfg.App.Mode == "debug" {
			ctx.Next()
			return
		}

		// 校验跨域请求
		if _, ok := m.cfg.AllowedOriginsMap[origin]; !ok {
			// 拒绝跨域
			m.log.Warn(ctx, "跨域", logger.StringAny("origin", origin))
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "跨域来源不允许",
			})
			return
		}

		// 处理请求
		ctx.Next()
	}
}

// setHeaders 设置允许跨域请求的响应头
func setHeaders(ctx *gin.Context, origin string) {
	ctx.Header("Access-Control-Allow-Origin", origin)
	ctx.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Token,X-User-Id")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type,New-Token,New-Expires-At")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Max-Age", "86400")
}
