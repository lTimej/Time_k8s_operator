package myredis

import (
	"context"
	"time"

	"Time_k8s_operator/conf"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         conf.RedisConfig.Addr,
		Password:     conf.RedisConfig.Password,
		DB:           int(conf.RedisConfig.DB),
		PoolSize:     int(conf.RedisConfig.PoolSize),     // 连接池最大socket连接数
		MinIdleConns: int(conf.RedisConfig.MinIdleConns), // 最少连接维持数
	})

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := RedisClient.Ping(timeoutCtx).Result()
	if err != nil {
		return err
	}
	logrus.Println("config redis inited... ...")
	return nil
}

func CloseRedisConn() {
	RedisClient.Close()
}
