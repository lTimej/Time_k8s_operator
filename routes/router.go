package routes

import (
	"Time_k8s_operator/conf"

	"github.com/gin-gonic/gin"
)

func NewRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	var router *gin.Engine
	//开发环境
	if conf.ServerConfig.Mode == "Dev" {
		router = gin.Default()
	} else { //线上环境
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	}
	if len(middlewares) > 0 {
		// router.Use(gin.RecoveryWithWriter(logger.Outer()))
		router.Use(middlewares...)
	}
	return router
}
