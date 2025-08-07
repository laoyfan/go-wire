package middleware

import (
	"fmt"
	"go-wire/config"
	"go-wire/constant"
	"go-wire/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorMiddleware struct {
	cfg *config.Config
	log logger.Logger
}

func NewErrorMiddleware(cfg *config.Config, log logger.Logger) *ErrorMiddleware {
	return &ErrorMiddleware{cfg: cfg, log: log}
}

func (m *ErrorMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				m.handlePanic(ctx, m.cfg.App.Name, r)
			}
		}()
		ctx.Next()
	}
}

func (m *ErrorMiddleware) handlePanic(ctx *gin.Context, tag string, r any) {
	if be, ok := r.(constant.ErrorResponse); ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": be.Code,
			"msg":  be.Msg,
		})
		return
	}

	var brokenPipe bool

	if ne, ok := r.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			errStr := strings.ToLower(se.Error())
			if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
				brokenPipe = true
			}
		}
	}

	// 请求快照
	req, _ := httputil.DumpRequest(ctx.Request, false)
	requestInfo := strings.ReplaceAll(string(req), "\r\n", " ")

	// 错误信息
	errMsg := fmt.Sprintf("%v", r)
	stack := m.formatStackTrace(tag)

	if brokenPipe {
		m.log.Warn(ctx, "连接已断开",
			logger.StringAny("type", "broken_pipe"),
			logger.StringAny("error", errMsg),
			logger.StringAny("request", requestInfo),
		)
		ctx.Abort()
		return
	}

	// 是否是业务 error，可以加一个类型断言判断业务 error（可选）
	m.log.Error(ctx, "服务异常捕获",
		logger.StringAny("type", "panic"),
		logger.StringAny("error", errMsg),
		logger.StringAny("request", requestInfo),
		logger.StringAny("stack", stack),
	)

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"code": constant.ERROR,
		"msg":  "服务器开小差，请稍后再试",
	})
}

func (m *ErrorMiddleware) formatStackTrace(tag string) []string {
	stack := debug.Stack()
	lines := strings.Split(string(stack), "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "goroutine") {
			continue
		}
		if len(trimmedLine) > 0 && strings.Contains(trimmedLine, tag) {
			result = append(result, trimmedLine)
		}
	}
	return result
}
