//   File Name:  vocl.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/8 16:35
//    Change Activity:

package volc

import (
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	rtcbase "github.com/volcengine/volc-sdk-golang/base"
	"github.com/volcengine/volc-sdk-golang/service/rtc"
)

var VoclClient *rtc.Rtc

var VConfig Config

func InitWithViper(config viper.Viper) {
	utils.ReadViperConfig(config, "vocl", &VConfig)
	NewClient(VConfig)
}

func NewClient(config Config) {
	VoclClient = rtc.NewInstanceWithRegion(config.Region)
	VoclClient.SetCredential(rtcbase.Credentials{
		AccessKeyID:     config.AccessKey,
		SecretAccessKey: config.SecretKey,
	})
}
