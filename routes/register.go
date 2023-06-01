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
		resource_group.GET("/cpu", Response(resource_controller.GetCpu))
		resource_group.GET("/memory", Response(resource_controller.GetMemory))
		resource_group.GET("/disk", Response(resource_controller.GetDisk))
		resource_group.GET("/network", Response(resource_controller.GetNetwork))
	}

	space_template_group := router.Group("/space/template")
	space_controller := controller.NewCodeServiceController()
	{
		space_template_group.GET("/get", Response(space_controller.GetTemplateSpace))
		space_template_group.POST("/create", Response(space_controller.CreateTemplateSpace))
		space_template_group.PUT("/update/:st_id", Response(space_controller.EditTemplateSpace))
		space_template_group.DELETE("/delete/:st_id", Response(space_controller.DeleteTemplateSpace))
	}

	space_group := router.Group("/space")
	{
		space_group.GET("/get", Response(space_controller.GetSpace))
		space_group.POST("/create", Response(space_controller.CreateSpace))
		space_group.POST("/create/run", Response(space_controller.CreateSpaceAndRun))
		space_group.PUT("/stop", Response(space_controller.StopSpace))
		space_group.PUT("/start", Response(space_controller.StartSpace))
		space_group.DELETE("/delete", Response(space_controller.DeleteSpace))
	}
}
