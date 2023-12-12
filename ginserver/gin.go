/*
File Name:  gin.go
Description:
Author:      Chenghu
Date:       2023/9/27 16:35
Change Activity:
*/
package ginserver

import (
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
)

var GinConfig ServerConfig

type ServerConfig struct {
	Name                string `json:"name"`
	Addr                string `json:"addr"`
	HideServerMiddleLog bool   `json:"hideServerMiddleLog"` // 是否隐藏内置中间件的 http 日志
}

// 使用 viper读取的配置初始化
func InitGinWithViperConfig(config viper.Viper) {
	//GinConfig = readServerConfig(config)
	utils.ReadViperConfig(config, "gin", &GinConfig)
}
