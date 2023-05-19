package main

import (
	"Time_k8s_operator/conf"
	"Time_k8s_operator/pkg/logger"
	"Time_k8s_operator/routes"
	"fmt"
)

func main() {

	if err := conf.InitConf(); err != nil {
		panic(fmt.Errorf("load conf failed, reason:%s", err.Error()))
	}

	if err := logger.InitLogger(); err != nil {
		panic(fmt.Errorf("load logger failed, reason:%s", err.Error()))
	}

	// 创建gin路由
	engine := routes.NewRouter()
	// 注册路由
	routes.Register(engine)
	engine.Run(":8998")
}
