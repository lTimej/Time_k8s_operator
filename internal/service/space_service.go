package service

import (
	"Time_k8s_operator/internal/caches"
	"Time_k8s_operator/internal/dao"
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/internal/rpc"
	"Time_k8s_operator/pb"
	"Time_k8s_operator/pkg/logger"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

const (
	CloudCodeNamespace = "k8s-test"
	DefaultPodPort     = 9999
	MaxSpaceCount      = 20
)

var (
	ErrReqParamInvalid    = errors.New("参数错误")
	ErrNameDuplicate      = errors.New("空间名称重复")
	ErrReachMaxSpaceCount = errors.New("达到最大空间数量")
	ErrSpaceCreate        = errors.New("空间创建失败")
	ErrSpaceStart         = errors.New("空间启动失败")
	ErrSpaceAlreadyExist  = errors.New("空间已经存在")
	ErrSpaceNotFound      = errors.New("空间不存在")
	ErrResourceExhausted  = errors.New("资源不足")
)

type CodeService struct {
	logger        *logrus.Logger
	rpc           pb.ServiceClient
	templateCache *caches.TemplateCache
	specCache     *caches.SpaceCache
}

func generateSID() string {
	return bson.NewObjectId().Hex()
}

func NewCodeService() *CodeService {
	conn := rpc.GrpcClient("space-code")
	return &CodeService{
		logger:        logger.Logger(),
		rpc:           pb.NewServiceClient(conn),
		templateCache: caches.NewCacheFactory().TemplateCache(),
		specCache:     caches.NewCacheFactory().SpaceSpecCache(),
	}
}

func (cs *CodeService) CreateSpace(req model.SpaceCreateOption) (*model.Space, error) {
	count := dao.FindCountByUserId(req.UserId)
	if count > MaxSpaceCount {
		return nil, ErrReachMaxSpaceCount
	}

	if !dao.FindOneByUserIdAndName(req.UserId, req.Name) {
		return nil, ErrNameDuplicate
	}

	st := cs.templateCache.GetSpaceTemplate(req.TemplateId)
	if st == nil {
		cs.logger.Errorf("参数错误:%v", err)
		return nil, ErrReqParamInvalid
	}

	ss := cs.specCache.GetSpaceSpec(req.SpaceSpecId)
	if ss == nil {
		cs.logger.Errorf("参数错误:%v", err)
		return nil, ErrReqParamInvalid
	}
	now := time.Now()
	space := &model.Space{
		UserId:     req.UserId,
		TemplateId: st.Id,
		SpecId:     ss.Id,
		Spec:       *ss,
		Name:       req.Name,
		Status:     model.SpaceStatusUncreated,
		CreateTime: now,
		DeleteTime: now,
		StopTime:   now,
		TotalTime:  0,
		Sid:        generateSID(),
	}
	id, err := dao.InsertSpace(space)
	if err != nil {
		cs.logger.Errorf("创建空间失败:%v", err)
		return nil, ErrSpaceCreate
	}
	space.Id = id
	return space, nil
}
