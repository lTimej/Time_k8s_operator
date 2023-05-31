package controller

import (
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/internal/service"
	"Time_k8s_operator/pkg/code"
	"Time_k8s_operator/pkg/httpResp"
	"Time_k8s_operator/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CodeServiceController struct {
	logger      *logrus.Logger
	codeService *service.CodeService
}

func NewCodeServiceController() *CodeServiceController {
	return &CodeServiceController{
		logger:      logger.Logger(),
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
	csc.logger.Debug(reqInfo)

	// 参数验证
	// get1, exist1 := c.Get("id")
	// _, exist2 := c.Get("username")
	// if !exist1 || !exist2 {
	// 	fmt.Println("111111111111")
	// 	c.Status(http.StatusBadRequest)
	// 	return nil
	// }
	// id, ok := get1.(uint32)
	// if !ok || id != reqInfo.UserId {
	// 	fmt.Println(err, "11111133111111")
	// 	c.Status(http.StatusBadRequest)
	// 	return nil
	// }
	space, err := csc.codeService.CreateSpace(reqInfo)
	fmt.Println(err, "2222222222222222")
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

func (csc *CodeServiceController) CreateSpaceAndRun(c *gin.Context) *httpResp.Response {
	//获取参数
	var reqInfo model.SpaceCreateOption
	err := c.ShouldBind(&reqInfo)
	if err != nil {
		c.Status(http.StatusBadRequest)
		// return nil
	}
	//参数校验
	// csc.logger.Debug(reqInfo)
	space, err := csc.codeService.CreateSpaceAndRun(reqInfo)
	fmt.Println(err, "2222222222222222")
	switch err {
	case service.ErrNameDuplicate:
		return httpResp.ResponseOk(code.SpaceCreateNameDuplicate, nil)
	case service.ErrReachMaxSpaceCount:
		return httpResp.ResponseOk(code.SpaceCreateReachMaxCount, nil)
	case service.ErrSpaceCreate:
		return httpResp.ResponseOk(code.SpaceCreateFailed, nil)
	case service.ErrSpaceStart:
		return httpResp.ResponseOk(code.SpaceStartFailed, nil)
	case service.ErrOtherSpaceIsRunning:
		return httpResp.ResponseOk(code.SpaceOtherSpaceIsRunning, nil)
	case service.ErrReqParamInvalid:
		c.Status(http.StatusBadRequest)
		return nil
	case service.ErrSpaceAlreadyExist:
		return httpResp.ResponseOk(code.SpaceAlreadyExist, nil)
	case service.ErrResourceExhausted:
		return httpResp.ResponseOk(code.ResourceExhausted, nil)
	}

	if err != nil {
		return httpResp.ResponseOk(code.SpaceCreateFailed, nil)
	}

	return httpResp.ResponseOk(code.SpaceCreateSuccess, space)
}

func (csc *CodeServiceController) StopSpace(c *gin.Context) *httpResp.Response {
	var req struct {
		Id  uint32 `json:"id"`
		Sid string `json:"sid"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		csc.logger.Warningf("bind param error:%v", err)
		c.Status(http.StatusBadRequest)
		return nil
	}
	uid := "1"
	err = csc.codeService.StopSpace(req.Id, uid, req.Sid)
	if err != nil {
		if err == service.ErrWorkSpaceIsNotRunning {
			return httpResp.ResponseOk(code.SpaceStopIsNotRunning, nil)
		}
		return httpResp.ResponseOk(code.SpaceStopFailed, nil)
	}
	return httpResp.ResponseOk(code.SpaceStopSuccess, nil)
}

func (csc *CodeServiceController) StartSpace(c *gin.Context) *httpResp.Response {
	var req struct {
		Id uint32 `json:"id"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		// c.logger.Warnf("获取参数失败:%v", err)
		c.Status(http.StatusBadRequest)
		return nil
	}
	space, err := csc.codeService.StartSpace(req.Id, 1, "1")
	switch err {
	case service.ErrWorkSpaceNotExist:
		return httpResp.ResponseOk(code.SpaceStartNotExist, nil)
	case service.ErrSpaceStart:
		return httpResp.ResponseOk(code.SpaceStartFailed, nil)
	case service.ErrOtherSpaceIsRunning:
		return httpResp.ResponseOk(code.SpaceOtherSpaceIsRunning, nil)
	case service.ErrSpaceNotFound:
		return httpResp.ResponseOk(code.SpaceNotFound, nil)
	}

	if err != nil {
		return httpResp.ResponseOk(code.SpaceStartFailed, nil)
	}

	return httpResp.ResponseOk(code.SpaceStartSuccess, space)
}

func (csc *CodeServiceController) DeleteSpace(c *gin.Context) *httpResp.Response {
	var req struct {
		Id uint32 `json:"id"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		csc.logger.Warnf("获取参数失败:%v", err)
		c.Status(http.StatusBadRequest)
		return nil
	}
	csc.logger.Debug("id:", req.Id)
	err = csc.codeService.DeleteSpace(req.Id)
	if err != nil {
		if err == service.ErrOtherSpaceIsRunning {
			return httpResp.ResponseOk(code.SpaceDeleteIsRunning, nil)
		}
		return httpResp.ResponseOk(code.SpaceDeleteFailed, nil)
	}
	return httpResp.ResponseOk(code.SpaceDeleteSuccess, nil)
}
