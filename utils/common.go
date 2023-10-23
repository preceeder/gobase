//   File Name:  common.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/23 15:22
//    Change Activity:

package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

/** 将viper 配置读渠道 target 结构体中
 */
func ReadViperConfig(v viper.Viper, key string, target any) {
	keyConfig := v.Sub(key)
	if keyConfig == nil {
		fmt.Printf("%s config is nil\n", key)
		os.Exit(1)
	}
	err := keyConfig.Unmarshal(target)
	if err != nil {
		fmt.Printf("%s config read error: %s \n", key, err.Error())
		os.Exit(1)
	}
	return
}
