package util

import (
	"errors"

	"github.com/MetaverseTopDJ/Scaffold/model"
)

type RedisConfig struct {
	Name string `json:"name"`
	Db   int    `json:"db"`
}

// GetRedisConfig 获取 Redis 配置
func GetRedisConfig(path, name string) (conf *RedisConfig, err error) {
	list := &model.RedisMapConfig{}
	err = ParseConfig(path, list)
	for n, v := range list.List {
		if n == name {
			conf.Name = n
			conf.Db = v.Db
		}
	}
	if err != nil {
		err = errors.New("ReadRedisConfigError")
	}
	return
}
