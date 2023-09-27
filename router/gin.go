/*
File Name:  gin.py
Description:
Author:      Chenghu
Date:       2023/9/27 16:35
Change Activity:
*/
package router

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var GinConfig ServerConfig

type ServerConfig struct {
	Name string
	Addr string
}

// 使用 viper读取的配置初始化
func InitGinWithViperConfig(config viper.Viper) {
	GinConfig = readServerConfig(config)
}

func readServerConfig(config viper.Viper) ServerConfig {
	sc := config.Sub("server")
	if sc == nil {
		fmt.Printf("gin config is nil")
		os.Exit(1)
	}
	return ServerConfig{
		Name: sc.GetString("name"),
		Addr: sc.GetString("addr"),
	}
}
