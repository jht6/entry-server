package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() {
	env := GetServerEnv()

	file := fmt.Sprintf("./config/%s.yaml", env)
	viper.SetConfigFile(file)

	fmt.Printf("Using config file: %v\n", file)

	err := viper.ReadInConfig()
	if err != nil {
		panic("读取配置文件时发生异常，服务启动失败")
	}
}

// 读取配置，统一返回string，由使用处自行转换类型
func CfgGet(key string) string {
	return viper.GetString(key)
}
