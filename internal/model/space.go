package model

import "time"

// Space的Status
const (
	SpaceStatusDeleted   = iota //被删除的空间状态
	SpaceStatusAvailable        //创建空间之后pod启动之后状态
	SpaceStatusUncreated        //创建空间之后pod启动之前状态
)

const (
	RunningStatusStop    = iota //空间运行停止
	RunningStatusRunning        //空间正在运行
)

type TemplateKind struct {
	Id   uint32 `gorm:"type:bigint(20)"`
	Name string `gorm:"type:varchar(255);"`
}

func (tk *TemplateKind) TableName() string {
	return "template_kind"
}

// SpaceSpec 云空间的配置
type SpaceSpec struct {
	Id          uint32 `gorm:"type:bigint(20)"`
	CpuSpec     string `gorm:"type:varchar(255);"` // CPU规格
	MemSpec     string `gorm:"type:varchar(255);"` // 内存规格
	StorageSpec string `gorm:"type:varchar(255);"` // 存储规格
	Name        string `gorm:"type:varchar(255);"`
	Description string `gorm:"type:varchar(255);"`
}

func (ss *SpaceSpec) TableName() string {
	return "space_spec"
}

type Space struct {
	Id            uint32        `gorm:"type:bigint(20)"`
	UserId        uint32        `gorm:"type:bigint(20)"` // 所属用户的id
	TemplateId    uint32        `gorm:"type:bigint(20)"` // 模板的id
	SpecId        uint32        `gorm:"type:bigint(20)"` // 规格id
	Spec          SpaceSpec     `gorm:"ForeignKey:SpecId;AssociationForeignKey:ID"`
	Sid           string        `gorm:"type:varchar(255);"` // 工作空间Id，用于访问时的url中
	Name          string        `gorm:"type:varchar(255);"` // 名称
	Status        uint32        `gorm:"type:tinyint(1)"`    // 0 已删除  1 可用 2 未创建
	RunningStatus uint32        `gorm:"type:tinyint(1)"`    // 0 停止  1 正在运行
	CreateTime    time.Time     `gorm:"type:datetime; DEFAULT CURRENT_TIMESTAMP"`
	DeleteTime    time.Time     `gorm:"type:datetime; DEFAULT NULL"`
	StopTime      time.Time     `gorm:"type:datetime; DEFAULT NULL"`   // 停止时间
	TotalTime     time.Duration `gorm:"type:bigint(20); DEFAULT NULL"` // 总运行时间
	Environment   string        `gorm:"type:varchar(255);"`
	Avatar        string        `gorm:"type:varchar(255);"`
}

func (s *Space) TableName() string {
	return "space"
}

type SpaceTemplate struct {
	Id          uint32    `gorm:"type:bigint(20)"`
	KindId      uint32    `gorm:"type:bigint(20)"`    // 类别Id
	Name        string    `gorm:"type:varchar(255);"` // 空间模板名称
	Description string    `gorm:"type:varchar(255);"` // 描述
	Tags        string    `gorm:"type:varchar(255);"` // 标签，使用|隔开
	Image       string    `gorm:"type:varchar(255);"` // 镜像
	Status      uint32    `gorm:"type:tinyint(1)"`    // 0可用 1 已删除
	Avatar      string    `gorm:"type:varchar(255);"`
	CreateTime  time.Time `gorm:"type:datetime; DEFAULT NULL"`
	DeleteTime  time.Time `gorm:"type:datetime; DEFAULT NULL"`
}

func (st *SpaceTemplate) TableName() string {
	return "space_template"
}

type RunningSpace struct {
	Sid  string `gorm:"type:varchar(255);"`
	Host string `gorm:"type:varchar(255);"`
}

func (rs *RunningSpace) TableName() string {
	return "running_space"
}

type SpaceCreateOption struct {
	Name        string `json:"name"`
	TemplateId  uint32 `json:"template_id"`
	SpaceSpecId uint32 `json:"space_spec_id"`
	UserId      uint32 `json:"user_id"`
}

type SpaceTemplateCreateOption struct {
	KindId      uint32 `json:"kind_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags"`   //标签
	Image       string `json:"image"`  //镜像
	Avatar      string `json:"avatar"` //logo
}
