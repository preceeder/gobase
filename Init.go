/*
File Name:  Init.py
Description:   这里一定要初始化
Author:      Chenghu
Date:       2023/8/21 13:35
Change Activity:
*/
package gobase

import (
	"github.com/preceeder/gobase/aliyun/push"
	"github.com/preceeder/gobase/config"
	"github.com/preceeder/gobase/db/mysqlDb"
	"github.com/preceeder/gobase/db/redisDb"
	"github.com/preceeder/gobase/env"
	"github.com/preceeder/gobase/ginserver"
	"github.com/preceeder/gobase/grpcm"
	"github.com/preceeder/gobase/jigou"
	"github.com/preceeder/gobase/logs"
	"github.com/preceeder/gobase/rongyun"
	"github.com/preceeder/gobase/shumei"
	"github.com/spf13/viper"
)

type initOptional struct {
	withRedis  bool
	withMysql  bool
	withIm     bool
	withJigou  bool
	withRpc    bool
	withShumei bool
}

func WithGinOptional(c bool) func(*initOptional, viper.Viper) {
	return func(il *initOptional, config viper.Viper) {
		il.withRedis = c
		if c == true {
			ginserver.InitGinWithViperConfig(config)
		}
	}
}

func WithRedisOptional(c bool) func(*initOptional, viper.Viper) {
	return func(il *initOptional, config viper.Viper) {
		il.withRedis = c
		if c == true {
			redisDb.InitRedisWithViperConfig(config)
		}
	}
}

func WithMysqlOptional(c bool) func(*initOptional, viper.Viper) {
	return func(il *initOptional, config viper.Viper) {
		il.withMysql = c
		if c == true {
			mysqlDb.InitMysqlWithViperConfig(config)

		}
	}
}

func WithImOptional(c bool) func(*initOptional, viper.Viper) {
	return func(il *initOptional, config viper.Viper) {
		il.withIm = c
		if c == true {
			rongyun.InitWithViper(config)
		}
	}
}
func WithJigouOptional(c bool) func(*initOptional, viper.Viper) {
	return func(il *initOptional, config viper.Viper) {
		il.withJigou = c
		if c == true {
			jigou.InitWithViper(config)

		}
	}
}

func WithRpcOptional(c bool) func(*initOptional, viper.Viper) {
	return func(il *initOptional, config viper.Viper) {
		il.withRpc = c
		if c == true {
			grpcm.InitRpc(config)
		}
	}
}

func WithAliYunPushOptional(c bool) func(*initOptional, viper.Viper) {
	return func(il *initOptional, config viper.Viper) {
		il.withRpc = c
		if c == true {
			push.InitWithViper(config)
		}
	}
}

func WithShumeiOptional(c bool) func(*initOptional, viper.Viper) {
	return func(il *initOptional, config viper.Viper) {
		il.withShumei = c
		if c == true {
			shumei.InitShumeiWithViperConfig(config)
		}
	}
}

func Init(viperPath string, viperConfigName string, optional ...func(*initOptional, viper.Viper)) {
	cf := config.InitConfig(viperPath, viperConfigName)
	//初始化环境变量
	env.InitEnv(*cf.Viper)

	logs.InitLogWithViper(*cf.Viper)

	il := initOptional{}
	for _, v := range optional {
		v(&il, *cf.Viper)
	}

	return
}
