package middleware

import (
	"Time_k8s_operator/pkg/code"
	"Time_k8s_operator/pkg/httpResp"
	"Time_k8s_operator/pkg/utils/encrypt"
	"fmt"
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
		myclaim, err := encrypt.VerifyToken(token)
		switch err {
		case encrypt.ErrToKenInvalid:
			context.Abort()
			httpResp.HttpResp(context, code.ErrToKenInvalid, nil)
			return
		case encrypt.ErrTokenExpire:
			context.Abort()
			httpResp.HttpResp(context, code.ErrTokenExpire, nil)
			return
		}
		username := myclaim.Username
		id := myclaim.Id
		uid := myclaim.Uid
		fmt.Println(username, id, uid, "=========")
		context.Set("username", username)
		context.Set("id", id)
		context.Set("uid", uid)
		context.Next()
	}
}
