package middleware

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAuthMiddleware,
	NewCorsMiddleware,
	NewLimiter,
	NewLimiterMiddleware,
	NewTraceMiddleware,
	NewErrorMiddleware,
	NewLoggerMiddleware,
)
