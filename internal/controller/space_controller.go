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

type CodeServiceController struct {
	logger      *logrus.Logger
	codeService *service.CodeService
}

func NewCodeServiceController() *CodeServiceController {
	return &CodeServiceController{
		logger: logger.Logger(),

		codeService: service.NewCodeService(),
	}
}

func (csc *CodeServiceController) GetSpaceSpec(c *gin.Context) *httpResp.Response {
	space_specs := csc.codeService.GetSpaceSpec()
	return httpResp.ResponseOk(code.SpaceSpecGetSuccess, space_specs)
}

func (csc *CodeServiceController) CreateSpaceSpec(c *gin.Context) *httpResp.Response {
	//获取参数
	var reqInfo model.SpaceSpecCreateOption
	err := c.ShouldBind(&reqInfo)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return nil
	}
	//参数校验
	csc.logger.Debug(reqInfo)
	space_spec, err := csc.codeService.CreateSpaceSpec(reqInfo)
	switch err {
	case service.ErrSpaceSpecCreate:
		return httpResp.ResponseOk(code.ErrSpaceSpecCreate, nil)
	case service.ErrReqParamInvalid:
		c.Status(http.StatusBadRequest)
		return nil
	}
	if err != nil {
		return httpResp.ResponseOk(code.ErrSpaceSpecCreate, nil)
	}
	return httpResp.ResponseOk(code.SpaceSpecCreateSuccess, space_spec)
}

func (csc *CodeServiceController) GetTemplateKind(c *gin.Context) *httpResp.Response {
	template_kinds := csc.codeService.GetTemplateKind()
	return httpResp.ResponseOk(code.TemplateKindGetSuccess, template_kinds)
}

func (csc *CodeServiceController) GetTemplateSpace(c *gin.Context) *httpResp.Response {
	space_templates := csc.codeService.GetTemplateSpace()
	return httpResp.ResponseOk(code.SpaceTemplateGetSuccess, space_templates)
}

func (csc *CodeServiceController) CreateTemplateSpace(c *gin.Context) *httpResp.Response {
	//获取参数
	var reqInfo model.SpaceTemplateCreateOption
	err := c.ShouldBind(&reqInfo)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return nil
	}
	//参数校验
	csc.logger.Debug(reqInfo)
	space_template, err := csc.codeService.CreateTemplateSpace(reqInfo)
	switch err {
	case service.ErrSpaceTemplateAlreadyExist:
		return httpResp.ResponseOk(code.SpaceTemplateCreateNameDuplicate, nil)
	case service.ErrSpaceTemplateCreate:
		return httpResp.ResponseOk(code.ErrSpaceTemplateCreate, nil)
	case service.ErrReqParamInvalid:
		c.Status(http.StatusBadRequest)
		return nil
	}
	if err != nil {
		return httpResp.ResponseOk(code.ErrSpaceTemplateCreate, nil)
	}
	return httpResp.ResponseOk(code.SpaceTemplateCreateSuccess, space_template)
}

func (csc *CodeServiceController) EditTemplateSpace(c *gin.Context) *httpResp.Response {
	//获取参数
	var reqInfo model.SpaceTemplateCreateOption
	err := c.ShouldBind(&reqInfo)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return nil
	}
	st_id := c.Param("st_id")
	//参数校验
	csc.logger.Debug(reqInfo)
	space_template, err := csc.codeService.EditTemplateSpace(reqInfo, st_id)
	switch err {
	case service.ErrSpaceTemplateNotExist:
		return httpResp.ResponseOk(code.SpaceTemplateNotExist, nil)
	case service.ErrSpaceTemplateUpdate:
		return httpResp.ResponseOk(code.ErrSpaceTemplateUpdateFailed, nil)
	case service.ErrReqParamInvalid:
		c.Status(http.StatusBadRequest)
		return nil
	}
	if err != nil {
		return httpResp.ResponseOk(code.ErrSpaceTemplateUpdateFailed, nil)
	}
	return httpResp.ResponseOk(code.SpaceTemplateUpdateSuccess, space_template)
}
func (csc *CodeServiceController) DeleteTemplateSpace(c *gin.Context) *httpResp.Response {
	//获取参数
	st_id := c.Param("st_id")
	//参数校验
	csc.logger.Debug(st_id)
	err := csc.codeService.DeleteTemplateSpace(st_id)
	switch err {
	case service.ErrSpaceTemplateNotExist:
		return httpResp.ResponseOk(code.SpaceTemplateNotExist, nil)
	case service.ErrSpaceTemplateUpdate:
		return httpResp.ResponseOk(code.ErrSpaceTemplateDelete, nil)
	case service.ErrReqParamInvalid:
		c.Status(http.StatusBadRequest)
		return nil
	}
	if err != nil {
		return httpResp.ResponseOk(code.ErrSpaceTemplateDelete, nil)
	}
	return httpResp.ResponseOk(code.SpaceTemplateDeleteSuccess, nil)
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

func (csc *CodeServiceController) CreateSpaceAndRun(c *gin.Context) *httpResp.Response {
	//获取参数
	var reqInfo model.SpaceCreateOption
	err := c.ShouldBind(&reqInfo)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return nil
	}
	uidi, ok := c.Get("uid")
	if !ok {
		return httpResp.ResponseOk(code.UserNotLogin, nil)
	}
	uid := uidi.(string)
	//参数校验
	csc.logger.Debug(reqInfo)
	space, err := csc.codeService.CreateSpaceAndRun(reqInfo, uid)
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
	uidi, ok := c.Get("uid")
	if !ok {
		return httpResp.ResponseOk(code.UserNotLogin, nil)
	}
	uid := uidi.(string)
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
		csc.logger.Warnf("获取参数失败:%v", err)
		c.Status(http.StatusBadRequest)
		return nil
	}
	uidi, ok := c.Get("uid")
	if !ok {
		return httpResp.ResponseOk(code.UserNotLogin, nil)
	}
	uid := uidi.(string)
	space, err := csc.codeService.StartSpace(req.Id, 1, uid)
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
	uidi, ok := c.Get("uid")
	if !ok {
		return httpResp.ResponseOk(code.UserNotLogin, nil)
	}
	uid := uidi.(string)
	err = csc.codeService.DeleteSpace(req.Id, uid)
	if err != nil {
		if err == service.ErrOtherSpaceIsRunning {
			return httpResp.ResponseOk(code.SpaceDeleteIsRunning, nil)
		}
		return httpResp.ResponseOk(code.SpaceDeleteFailed, nil)
	}
	return httpResp.ResponseOk(code.SpaceDeleteSuccess, nil)
}

func (csc *CodeServiceController) GetSpace(c *gin.Context) *httpResp.Response {
	user_idi, ok := c.Get("id")
	if !ok {
		return httpResp.ResponseOk(code.UserNotLogin, nil)
	}
	user_id := user_idi.(uint32)
	spaces := csc.codeService.GetSpace(user_id)
	return httpResp.ResponseOk(code.SpaceGetSuccess, spaces)
}
