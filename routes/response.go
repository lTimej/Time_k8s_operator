package routes

import (
	"Time_k8s_operator/pkg/httpResp"

	"github.com/gin-gonic/gin"
)

type Handler func(*gin.Context) *httpResp.Response

func Response(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := h(c)
		if r != nil {
			c.JSON(r.HttpStatus, &r.Result)
		}
		httpResp.PutResponse(r)
	}
}
