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
	"log/slog"
	"time"
)

type VoclClientA struct {
	Client *rtc.Rtc
	Config Config
}

var VoclClient VoclClientA

func InitWithViper(config viper.Viper) {
	VoclClient = VoclClientA{}
	utils.ReadViperConfig(config, "vocl", &VoclClient.Config)
	NewClient(VoclClient.Config)
}

func NewClient(config Config) {
	VoclClient.Client = rtc.NewInstanceWithRegion(config.Region)
	VoclClient.Client.SetCredential(rtcbase.Credentials{
		AccessKeyID:     config.AccessKey,
		SecretAccessKey: config.SecretKey,
	})
}

func (v VoclClientA) GetToken(roomId, userId string) (string, error) {
	token := New(v.Config.AppId, v.Config.AppKey, roomId, userId)
	token.ExpireTime(time.Now().Add(time.Hour * 24))
	token.AddPrivilege(PrivSubscribeStream, time.Time{})
	token.AddPrivilege(PrivPublishStream, time.Now().Add(time.Minute))
	s, err := token.Serialize()
	if err != nil {
		slog.Error("生成token失败", "error", err.Error())
		return "", err
	}
	return s, nil
}
