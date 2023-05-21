package service

import (
	"Time_k8s_operator/internal/dao"
	"Time_k8s_operator/internal/dao/db"
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/pkg/logger"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrEmailCodeIncorrect = errors.New("email code incorrect")
	ErrEmailAlreadyInUse  = errors.New("this email had been registered")
	DbErr                 = errors.New("db error")
	ErrUserOrPasswordNOtIncorrect = errors.New("用户名或密码错误")
	ErrUserNotPresent = errors.New("用户不存在")
)

type UserService struct {
	logger *logrus.Logger
	// dao    *dao.UserDao
}

func NewUserService() *UserService {
	return &UserService{
		logger: logger.Logger(),
		//		dao:    dao.NewUserDao(),
	}
}

func (u *UserService) Login(username,password string) error{
	if dao.FindOneUserByUsername(username){
		u.logger.Infof("用户不存在")
		return ErrUserNotPresent
	}
	if dao.FindOneUserByUsernameAndPassword(username,password){
		u.logger.Infof("用户名或密码错误")
		return ErrUserOrPasswordNOtIncorrect
	}

	return nil
}

func (u *UserService) Register(register_info model.RegisterInfo) error {
	if dao.FindOneUserByEmail(register_info.Email) {
		u.logger.Infof("邮箱已存在!")
		return ErrEmailAlreadyInUse
	}
	now := time.Now()
	user := &model.User{
		Uid:        bson.NewObjectId().Hex(),
		Username:   register_info.Username,
		Password:   register_info.Password,
		Nickname:   register_info.Nickname,
		Email:      register_info.Email,
		CreateTime: now,
		DeleteTime: now,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		return DbErr
	}
	return nil
}
