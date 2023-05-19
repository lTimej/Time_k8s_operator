package service

import (
	"Time_k8s_operator/pkg/logger"

	"github.com/sirupsen/logrus"
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
