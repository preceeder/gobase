//   File Name:  common.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/23 15:22
//    Change Activity:

package utils

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
)

/** 将viper 配置读渠道 target 结构体中
 */
func ReadViperConfig(v viper.Viper, key string, target any) {
	err := v.UnmarshalKey(key, target, func(ms *mapstructure.DecoderConfig) { ms.TagName = "json" })
	if err != nil {
		fmt.Printf("load %s config error: %s\n", key, err.Error())
		os.Exit(1)
	}

	return
}
