package service

import (
	"Time_k8s_operator/internal/caches"
	"Time_k8s_operator/internal/dao"
	"Time_k8s_operator/internal/dao/myredis"
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/internal/rpc"
	"Time_k8s_operator/pb"
	"Time_k8s_operator/pkg/logger"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
)

const (
	CodeNamespace  = "k8s-test"
	DefaultPodPort = 9999
	MaxSpaceCount  = 20
)

var (
	ErrSpaceTemplateCreate       = errors.New("空间模板创建失败")
	ErrSpaceTemplateUpdate       = errors.New("空间模板修改失败")
	ErrSpaceTemplateDelete       = errors.New("空间模板删除失败")
	ErrReqParamInvalid           = errors.New("参数错误")
	ErrNameDuplicate             = errors.New("空间名称重复")
	ErrReachMaxSpaceCount        = errors.New("达到最大空间数量")
	ErrSpaceCreate               = errors.New("空间创建失败")
	ErrSpaceStart                = errors.New("空间启动失败")
	ErrSpaceAlreadyExist         = errors.New("空间已经存在")
	ErrSpaceTemplateAlreadyExist = errors.New("空间模板已经存在")
	ErrSpaceTemplateNotExist     = errors.New("空间模板不存在")
	ErrSpaceNotFound             = errors.New("空间不存在")
	ErrResourceExhausted         = errors.New("资源不足")
	ErrOtherSpaceIsRunning       = errors.New("其他空间正在运行")
	ErrWorkSpaceIsNotRunning     = errors.New("工作空间没有运行")
	ErrWorkSpaceNotExist         = errors.New("工作空间不存在")
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
	factory := caches.NewCacheFactory()
	return &CodeService{
		logger:        logger.Logger(),
		rpc:           pb.NewServiceClient(conn),
		templateCache: factory.TemplateCache(),
		specCache:     factory.SpaceSpecCache(),
	}
}

func (cs *CodeService) GetTemplateKind() []*model.TemplateKind {
	template_kinds := cs.templateCache.GetAllKind()
	return template_kinds
}

func (cs *CodeService) GetTemplateSpace() []*model.SpaceTemplate {
	space_templates := dao.FindAllTemplate()
	return space_templates
}

func (cs *CodeService) CreateTemplateSpace(req model.SpaceTemplateCreateOption) (*model.SpaceTemplate, error) {
	_, ok := dao.FindOneTemplateByName(req.Name)
	if !ok {
		return nil, ErrSpaceTemplateAlreadyExist
	}
	now := time.Now()
	space_template := &model.SpaceTemplate{
		KindId:      req.KindId,
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
		Image:       req.Image,
		Avatar:      req.Avatar,
		Status:      1,
		CreateTime:  now,
		DeleteTime:  now,
	}
	st_id, err := dao.InsertSpaceTemplate(space_template)
	if err != nil {
		cs.logger.Errorf("创建空间模板失败:%v", err)
		return nil, ErrSpaceTemplateCreate
	}
	space_template.Id = st_id
	return space_template, nil
}

func (cs *CodeService) EditTemplateSpace(req model.SpaceTemplateCreateOption, st_id string) error {
	_, ok := dao.FindOneTemplateById(st_id)
	if ok {
		return ErrSpaceTemplateNotExist
	}
	err := dao.UpdateSpaceTemplate(req, st_id)
	if err != nil {
		cs.logger.Errorf("修改空间模板失败:%v", err)
		return ErrSpaceTemplateUpdate
	}
	return nil
}

func (cs *CodeService) DeleteTemplateSpace(st_id string) error {
	_, ok := dao.FindOneTemplateById(st_id)
	if ok {
		return ErrSpaceTemplateNotExist
	}
	err := dao.DeleteSpaceTemplate(st_id)
	if err != nil {
		cs.logger.Errorf("删除空间模板失败:%v", err)
		return ErrSpaceTemplateDelete
	}
	return nil
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
		UserId:        req.UserId,
		TemplateId:    st.Id,
		SpecId:        ss.Id,
		Spec:          *ss,
		Name:          req.Name,
		Status:        model.SpaceStatusUncreated,
		RunningStatus: model.RunningStatusStop,
		CreateTime:    now,
		DeleteTime:    now,
		StopTime:      now,
		TotalTime:     0,
		Sid:           generateSID(),
	}
	id, err := dao.InsertSpace(space)
	if err != nil {
		cs.logger.Errorf("创建空间失败:%v", err)
		return nil, ErrSpaceCreate
	}
	space.Id = id
	return space, nil
}

func (cs *CodeService) CreateSpaceAndRun(req model.SpaceCreateOption, uid string) (*model.Space, error) {
	//判断是否有其他空间正在运行
	isRunning, err := myredis.IsRunningSpace(uid)
	if err != nil {
		return nil, ErrSpaceCreate
	}
	if isRunning {
		return nil, ErrOtherSpaceIsRunning
	}
	space, err := cs.CreateSpace(req)
	if err != nil {
		return nil, err
	}
	return cs.runSpace(space, uid, cs.rpc.CreateSpace)
}

type StartFunc func(ctx context.Context, in *pb.WorkspaceInfo, opts ...grpc.CallOption) (*pb.WorkspaceRunningInfo, error)

func (cs *CodeService) runSpace(space *model.Space, uid string, startFunc StartFunc) (*model.Space, error) {
	space_template := cs.templateCache.GetSpaceTemplate(space.TemplateId)
	if space_template == nil {
		cs.logger.Warnf("获取模板缓存失败")
		return nil, ErrSpaceStart
	}
	pod_name := cs.genPodName(space.Sid, uid)
	fmt.Println(pod_name)
	pod := pb.WorkspaceInfo{
		Name:            pod_name,
		Namespace:       CodeNamespace,
		Image:           space_template.Image,
		Port:            DefaultPodPort,
		VolumeMountPath: "/user_data/",
		ResourceLimit: &pb.ResourceLimit{
			Cpu:     space.Spec.CpuSpec,
			Memory:  space.Spec.MemSpec,
			Storage: space.Spec.StorageSpec,
		},
	}
	var retErr error
Loop:
	for i := 0; i < 1; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		spaceInfo, err := startFunc(ctx, &pod)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				return nil, err
			}
			switch s.Code() {
			// 创建工作空间时,工作空间已存在,修改数据库中的status
			case codes.AlreadyExists:
				retErr = ErrSpaceAlreadyExist
				break Loop

			// 启动工作空间时,工作空间不存在
			case codes.NotFound:
				return nil, ErrSpaceNotFound

			// 资源耗尽,无法启动
			case codes.ResourceExhausted:
				return nil, ErrResourceExhausted
			case codes.Unknown:
				cs.logger.Errorf("rpc start space error:%v", err)
				return nil, ErrSpaceStart
			}
		}
		host := spaceInfo.Ip + ":" + strconv.Itoa(int(spaceInfo.Port))
		running_space := &model.RunningSpace{
			Host: host,
			Sid:  space.Sid,
		}
		err = myredis.RunningSpace(uid, running_space)
		if err != nil {
			cs.logger.Errorf("添加pod到redis失败:%v", err)
			return nil, ErrSpaceStart
		}
		space.RunningStatus = model.RunningStatusRunning
	}
	//修改数据库状态
	if space.Status == model.SpaceStatusUncreated {
		err := dao.UpdateSpaceStatusAndRunningStatus(space.Id, model.SpaceStatusAvailable, space.RunningStatus)
		if err != nil {
			cs.logger.Errorf("更新空间状态失败:%v", err)
			return nil, err
		}
	}
	if retErr != nil {
		return nil, retErr
	}
	return space, nil
}

func (cs *CodeService) StopSpace(space_id uint32, uid, sid string) error {
	//先检查运行空间是否存在
	isRunning, err := myredis.CheckRunningSpace(sid)
	if err != nil {
		cs.logger.Warnf("检查空间运行状态失败:%v", err)
		return err
	}
	if !isRunning {
		return ErrWorkSpaceIsNotRunning
	}
	//然后删出对应的pod
	pod_name := cs.genPodName(sid, uid)
	_, err = cs.rpc.StopSpace(context.Background(), &pb.QueryOption{
		Name:      pod_name,
		Namespace: CodeNamespace,
	})
	if err != nil {
		cs.logger.Warnf("pod删除失败err:%v", err)
		return err
	}
	//删除redis
	isRunning, err = myredis.DeleteRunningSpace(uid)
	if err != nil {
		cs.logger.Warnf("检查空间运行状态失败:%v", err)
		return err
	}
	if !isRunning {
		return ErrWorkSpaceIsNotRunning
	}
	//修改runningstatus状态
	err = dao.UpdateSpaceRunningStatus(space_id, model.RunningStatusStop)
	if err != nil {
		cs.logger.Warnf("空间状态修改失败err:%v", err)
	}
	return nil
}

func (cs *CodeService) StartSpace(id, user_id uint32, uid string) (*model.Space, error) {
	isRunning, err := myredis.IsRunningSpace(uid)
	if err != nil {
		return nil, ErrSpaceCreate
	}
	if isRunning {
		return nil, ErrOtherSpaceIsRunning
	}
	space, ok := dao.FindSpaceOneByIdAndUserId(id, user_id)
	if ok {
		cs.logger.Warnf("空间不存在")
		return nil, ErrWorkSpaceIsNotRunning
	}
	startFunc := cs.rpc.StartSpace
	switch space.Status {
	case model.SpaceStatusDeleted:
		return nil, ErrWorkSpaceNotExist
	case model.SpaceStatusUncreated:
		startFunc = cs.rpc.CreateSpace
		spec := cs.specCache.GetSpaceSpec(space.SpecId)
		if spec == nil {
			return nil, ErrSpaceStart
		}
		space.Spec = *spec
	}
	ret, err := cs.runSpace(&space, uid, startFunc)
	if err != nil {
		cs.logger.Warnf("启动工作空间失败:%v", err)
		return nil, err
	}
	//修改runningstatus状态
	err = dao.UpdateSpaceRunningStatus(id, model.RunningStatusRunning)
	if err != nil {
		cs.logger.Warnf("空间状态修改失败err:%v", err)
	}
	return ret, err
}

func (cs *CodeService) DeleteSpace(id uint32, uid string) error {
	space := dao.FindSpaceOneById(id)
	isRunning, err := myredis.CheckRunningSpace(space.Sid)
	if err != nil {
		cs.logger.Warnf("判断空间是否运行失败err:%v", err)
		return err
	}
	if isRunning {
		return ErrOtherSpaceIsRunning
	}
	pod_name := cs.genPodName(space.Sid, uid)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err = cs.rpc.DeleteSpace(ctx, &pb.QueryOption{
		Name:      pod_name,
		Namespace: CodeNamespace,
	})
	if err != nil {
		cs.logger.Warnf("删除工作空间失败:%v", err)
		return err
	}
	return dao.DeleteSpaceById(id)
}

func (cs *CodeService) genPodName(sid string, uid string) string {
	return strings.Join([]string{"ws", uid, sid}, "-")
}

func (cs *CodeService) GetSpace(user_id uint32) []*model.Space {
	spaces := dao.FindAllSpaceByUserId(user_id)
	for _, space := range spaces {
		spec := cs.specCache.GetSpaceSpec(space.SpecId)
		space.Spec = *spec
	}
	return spaces
}
