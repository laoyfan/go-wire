package middleware

import "github.com/gin-gonic/gin"

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
