package controller

import (
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/internal/service"
	"Time_k8s_operator/pkg/code"
	"Time_k8s_operator/pkg/httpResp"
	"Time_k8s_operator/pkg/logger"
	"net/http"

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
	//获取参数
	var login_info model.LoginInfo
	err := c.ShouldBind(&login_info)
	if err != nil {
		logrus.Println(err)
		c.Status(http.StatusBadRequest)
		return nil
	}
	username := login_info.Username
	password := login_info.Password
	//参数校验
	if username == "" {
		return httpResp.RepsonseNotOk("用户名不能为空")
	}
	if password == "" {
		return httpResp.RepsonseNotOk("密码不能为空")
	}
	user, err := u.userService.Login(username, password)
	switch err {
	case service.ErrUserNotPresent:
		return httpResp.ResponseOk(code.UserNameNotPresent, nil)
	case service.ErrUserOrPasswordNOtIncorrect:
		return httpResp.ResponseOk(code.UsernameOrPasswordErr, nil)
	case service.ErrGenToken:
		return httpResp.ResponseOk(code.GenTokenFailed, nil)
	case nil:
		return httpResp.ResponseOk(code.LoginSuccess, user)
	}
	return httpResp.ResponseOk(code.LoginFailed, nil)
}

func (u *UserController) Register(c *gin.Context) *httpResp.Response {
	// 获取参数
	var register_info model.RegisterInfo
	err := c.ShouldBind(&register_info)
	if err != nil {
		logrus.Println(err)
		c.Status(http.StatusBadRequest)
		return nil
	}
	// 参数校验
	if register_info.Password != register_info.RePassword {
		return httpResp.ResponseOk(code.PasswordNotunanimous, nil)
	}
	err = u.userService.Register(register_info)
	switch err {
	case service.ErrEmailCodeIncorrect:
		return httpResp.ResponseOk(code.UserEmailCodeIncorrect, nil)
	case service.ErrEmailAlreadyInUse:
		return httpResp.ResponseOk(code.UserEmailAlreadyInUse, nil)
	case service.DbErr:
		return httpResp.ResponseOk(code.DbErr, nil)
	case nil:
		return httpResp.ResponseOk(code.UserRegisterSuccess, nil)
	}
	logrus.Printf("注册失败:%v\n", err)
	return httpResp.ResponseOk(code.UserRegisterFailed, nil)
}
