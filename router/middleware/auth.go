package middleware

import (
	"go-wire/config"
	"go-wire/constant"
	"go-wire/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	cfg *config.Config
	log logger.Logger
}

func NewAuthMiddleware(cfg *config.Config, log logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg, log: log}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientID := ctx.Request.Header.Get("clientId")
		if clientID != m.cfg.App.ClientID {
			m.log.Warn(ctx, "无权限", logger.StringAny("clientId", clientID))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": constant.FORBIDDEN,
				"msg":  "无权限",
			})
			return
		}
		ctx.Next()
	}
}
