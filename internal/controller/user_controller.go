package controller

import (
	"Time_k8s_operator/internal/service"
	"Time_k8s_operator/pkg/code"
	"Time_k8s_operator/pkg/httpResp"
	"Time_k8s_operator/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	logger      *logrus.Logger
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		logger:      logger.Logger(),
		userService: service.NewUserService(),
	}
}

func (u *UserController) Login(c *gin.Context) *httpResp.Response {
	msg := "ok"
	return httpResp.ResponseOk(code.LoginSuccess, msg)
}

func (u *UserController) Register(c *gin.Context) *httpResp.Response {
	msg := "ok"
	return httpResp.ResponseOk(code.LoginSuccess, msg)
}
