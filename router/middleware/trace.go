package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TraceMiddleware struct{}

func NewTraceMiddleware() *TraceMiddleware {
	return &TraceMiddleware{}
}

func (m *TraceMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := uuid.New().String()
		ctx.Set("TraceID", traceID)
		ctx.Next()
	}
}
