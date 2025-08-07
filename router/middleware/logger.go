package middleware

import (
	"fmt"
	"go-wire/constant"
	"go-wire/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type LoggerMiddleware struct {
	log logger.Logger
}

func NewLoggerMiddleware(log logger.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{log: log}
}

func (m *LoggerMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()               // 记录请求开始时间
		path := ctx.Request.URL.Path      // 记录请求路径
		query := ctx.Request.URL.RawQuery // 记录请求Query参数

		ctx.Next()

		cost := time.Since(start)
		status := ctx.Writer.Status()

		// 构造日志信息
		layout := constant.LogLayout{
			Time:      start,
			Status:    status,
			Method:    ctx.Request.Method,
			Path:      path,
			Query:     query,
			IP:        ctx.ClientIP(),
			UserAgent: ctx.Request.UserAgent(),
			Error:     ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
			Cost:      cost.Seconds(),
			Source:    ctx.Request.Host,
		}

		msg := fmt.Sprintf("%d %v %s %s", status, cost.Milliseconds(), layout.Method, layout.Path)
		// 根据响应状态码决定记录日志级别
		if status >= 400 {
			m.log.Error(ctx, msg, logger.StringAny("log", layout))
		} else {
			m.log.Info(ctx, msg, logger.StringAny("log", layout))
		}
	}
}
