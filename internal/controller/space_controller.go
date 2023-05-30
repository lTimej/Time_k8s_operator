package controller

import (
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/internal/service"
	"Time_k8s_operator/pkg/code"
	"Time_k8s_operator/pkg/httpResp"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CodeServiceController struct {
	// logger      *logrus.Logger
	codeService *service.CodeService
}

func NewCodeServiceController() *CodeServiceController {
	return &CodeServiceController{
		// logger:      logger.Logger(),
		codeService: service.NewCodeService(),
	}
}

func (csc *CodeServiceController) CreateSpace(c *gin.Context) *httpResp.Response {
	//获取参数
	var reqInfo model.SpaceCreateOption
	err := c.ShouldBind(&reqInfo)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return nil
	}
	//参数校验
	// csc.logger.Debug(reqInfo)

	// 参数验证
	get1, exist1 := c.Get("id")
	_, exist2 := c.Get("username")
	if !exist1 || !exist2 {
		c.Status(http.StatusBadRequest)
		return nil
	}
	id, ok := get1.(uint32)
	if !ok || id != reqInfo.UserId {
		c.Status(http.StatusBadRequest)
		return nil
	}
	space, err := csc.codeService.CreateSpace(reqInfo)
	switch err {
	case service.ErrNameDuplicate:
		return httpResp.ResponseOk(code.SpaceCreateNameDuplicate, nil)
	case service.ErrReachMaxSpaceCount:
		return httpResp.ResponseOk(code.SpaceCreateReachMaxCount, nil)
	case service.ErrSpaceCreate:
		return httpResp.ResponseOk(code.SpaceCreateFailed, nil)
	case service.ErrReqParamInvalid:
		c.Status(http.StatusBadRequest)
		return nil
	}
	if err != nil {
		return httpResp.ResponseOk(code.SpaceCreateFailed, nil)
	}

	return httpResp.ResponseOk(code.SpaceCreateSuccess, space)
}
