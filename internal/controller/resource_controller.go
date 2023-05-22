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
