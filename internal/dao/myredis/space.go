package myredis

import (
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/pkg/utils"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

const (
	HostsHashKey = "hosts"
)

func RunningSpace(uid string, running_space *model.RunningSpace) error {
	data, err := json.Marshal(running_space)
	if err != nil {
		return err
	}
	res := utils.Bytes2String(data)
	ic := RedisClient.HSet(context.Background(), HostsHashKey, uid, res, running_space.Sid, running_space.Host)
	return ic.Err()
}

func IsRunningSpace(uid string) (bool, error) {
	res := RedisClient.HGet(context.Background(), HostsHashKey, uid)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		} else {
			return false, err
		}
	}
	if res.Val() == "" {
		return false, nil
	}
	return true, nil
}

func DeleteRunningSpace(uid string) (bool, error) {
	res := RedisClient.HGet(context.Background(), HostsHashKey, uid)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		} else {
			return false, err
		}
	}
	if res.Val() == "" {
		return false, nil
	}
	var space model.Space
	data := utils.String2Bytes(res.Val())
	err := json.Unmarshal(data, &space)
	if err != nil {
		return false, err
	}
	sid := space.Sid
	if err = RedisClient.HDel(context.Background(), HostsHashKey, uid, sid).Err(); err != nil {
		return true, err
	}

	return true, nil
}

func CheckRunningSpace(sid string) (bool, error) {
	res := RedisClient.HGet(context.Background(), HostsHashKey, sid)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		} else {
			return false, err
		}
	}
	if res.Val() == "" {
		return false, nil
	}
	return true, nil
}
