package controller

import (
	"Time_k8s_operator/internal/service"
	"Time_k8s_operator/pkg/code"
	"Time_k8s_operator/pkg/httpResp"
	"Time_k8s_operator/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ResourceController struct {
	logger          *logrus.Logger
	resourceService *service.ResourceService
}

func NewResourceController() *ResourceController {
	return &ResourceController{
		logger:          logger.Logger(),
		resourceService: service.NewResourceService(),
	}
}

func (rc ResourceController) GetResource(c *gin.Context) *httpResp.Response {
	data := rc.resourceService.GetResource()
	return httpResp.ResponseOk(code.GetResourceSuccess, data)
}

func (rc ResourceController) GetCpu(c *gin.Context) *httpResp.Response {
	data := rc.resourceService.GetCpu()
	return httpResp.ResponseOk(code.GetResourceSuccess, data)
}

func (rc ResourceController) GetMemory(c *gin.Context) *httpResp.Response {
	data := rc.resourceService.GetMemory()
	return httpResp.ResponseOk(code.GetResourceSuccess, data)
}

func (rc ResourceController) GetDisk(c *gin.Context) *httpResp.Response {
	data := rc.resourceService.GetDisk()
	return httpResp.ResponseOk(code.GetResourceSuccess, data)
}

func (rc ResourceController) GetNetwork(c *gin.Context) *httpResp.Response {
	data := rc.resourceService.GetNetwork()
	return httpResp.ResponseOk(code.GetResourceSuccess, data)
}
