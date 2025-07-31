package middleware

import "github.com/gin-gonic/gin"

func Limiter() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
