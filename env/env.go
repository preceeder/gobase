/*
File Name:  env.go
Description:
Author:      Chenghu
Date:       2023/8/29 15:56
Change Activity:
*/
package env

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

type Environment struct {
	Env      string `json:"env"` // product, test, local
	PreFix   string `json:"prefix"`
	GIN_MODE string `json:"GIN_MODE"`
}

func InitEnv(config viper.Viper) {
	//ev := ReadEnvConfig(config)
	ev := Environment{}
	utils.ReadViperConfig(config, "env", &ev)
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

func GetEnv(k string) string {
	e, ok := os.LookupEnv(k)
	if !ok {
		slog.Error("get env", "error", "variable is not defined", "key", k)
	}
	return e
}
