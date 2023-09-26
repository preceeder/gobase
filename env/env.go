/*
File Name:  env.py
Description:
Author:      Chenghu
Date:       2023/8/29 15:56
Change Activity:
*/
package env

import (
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

type Environment struct {
	Env    string `json:"env"` // product, test, local
	PreFix string `json:"prefix"`
}

func InitEnv(config viper.Viper) {
	ev := ReadEnvConfig(config)
	es, err := convertor.StructToMap(ev)
	if err != nil {
		panic("init env error: " + err.Error())
	}
	for k, v := range es {
		err := os.Setenv(k, v.(string))
		if err != nil {
			panic("init env error: " + err.Error())
		}
	}
}

func ReadEnvConfig(v viper.Viper) Environment {
	env := v.Sub("env")
	if env == nil {
		fmt.Printf("env config is nil")
		os.Exit(1)
	}
	ev := Environment{}
	err := env.Unmarshal(&ev)
	if err != nil {
		fmt.Printf("env config is nil")
		os.Exit(1)
	}
	return ev
}

func GetEnv(k string) string {
	e, ok := os.LookupEnv(k)
	if !ok {
		slog.Error("get env", "error", "variable is not defined", "key", k)
	}
	return e
}
