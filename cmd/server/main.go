package main

import (
	"Time_k8s_operator/conf"
	"Time_k8s_operator/internal/dao/db"
	"Time_k8s_operator/internal/dao/myredis"
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/pkg/logger"
	"Time_k8s_operator/routes"
	"fmt"
)

func main() {
	//初始化配置
	if err := conf.InitConf(); err != nil {
		panic(fmt.Errorf("load conf failed, reason:%s", err.Error()))
	}
	//初始化日志
	if err := logger.InitLogger(); err != nil {
		panic(fmt.Errorf("load logger failed, reason:%s", err.Error()))
	}
	//初始化mysql
	if err := db.InitMysql(); err != nil {
		panic(fmt.Errorf("load mysql failed, reason:%s", err.Error()))
	}
	//初始化redis
	if err := myredis.InitRedis(); err != nil {
		panic(fmt.Errorf("load redis failed, reason:%s", err.Error()))
	}
	//初始化数据库表
	if err := model.InitTable(); err != nil {
		panic(fmt.Errorf("load redis failed, reason:%s", err.Error()))
	}
	// 创建gin路由
	engine := routes.NewRouter()
	// 注册路由
	routes.Register(engine)
	engine.Run(":8998")
}
