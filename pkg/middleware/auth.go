package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.Request.Header.Get("Authorization")
		context.Next()
	}
}
