package middleware

import (
	"Time_k8s_operator/pkg/code"
	"Time_k8s_operator/pkg/httpResp"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.URL.Path == "/user/login" || context.Request.URL.Path == "/user/register" {
			context.Next()
			return
		}
		header := context.Request.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			context.Abort()
			httpResp.HttpResp(context, code.LoginNotAuth, nil)
			return
		}
		token := strings.Split(header, " ")[1]
		if token == "" || token == "null" {
			context.Abort()
			httpResp.HttpResp(context, code.LoginNotAuth, nil)
			return
		}
		context.Next()
	}
}
