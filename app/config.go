package app

import "strings"

// GetStringConfig 获取 string 格式的配置信息
func GetStringConfig(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return ""
	}
	confString := v.GetString(strings.Join(keys[1:], "."))
	return confString
}

// GetIntConfig 获取 int 格式的配置信息
func GetIntConfig(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetInt(strings.Join(keys[1:], "."))
	return conf
}
