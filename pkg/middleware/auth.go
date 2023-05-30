package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.URL.Path == "/user/login" {
			context.Next()
			return
		}
		header := context.Request.Header.Get("Authorization")
		fmt.Println(header, "====")
		if !strings.HasPrefix(header, "Bearer ") {
			return
		}
		token := strings.Split(header, " ")[1]
		if token == "" {
			return
		}
		context.Next()
	}
}
