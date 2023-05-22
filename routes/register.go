package routes

import (
	"Time_k8s_operator/internal/controller"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	user_group := router.Group("/user")
	user_controller := controller.NewUserController()
	{
		//登录
		user_group.POST("/login", Response(user_controller.Login))
		//注册
		user_group.POST("/register", Response(user_controller.Register))
	}

	resource_group := router.Group("/resource")
	resource_controller := controller.NewResourceController()
	{
		resource_group.GET("/getresource", Response(resource_controller.GetResource))
	}
}
