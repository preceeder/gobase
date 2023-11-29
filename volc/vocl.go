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
	"net/url"
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
	expt := time.Now().Add(time.Hour * 24)
	token.ExpireTime(expt)
	token.AddPrivilege(PrivSubscribeStream, expt)
	token.AddPrivilege(PrivPublishStream, expt)
	s, err := token.Serialize()
	if err != nil {
		slog.Error("生成token失败", "error", err.Error())
		return "", err
	}
	return s, nil
}

/**封禁房间, 用户 这里用来关闭房间
*appid  string 你的音视频应用的唯一标志
* roomId string 指定房间 ID
* userId string 希望封禁用户的 ID    只是关闭房间  传 ""
* ForbiddenInterval int 封禁时长，单位为秒  只是关闭房间  传 ""
 */
func (v VoclClientA) CloseRoom(roomId, appId, userId, forbiddenInterval string) error {
	form := url.Values{}
	if appId == "" {
		form.Set("AppId", v.Config.AppId)
	} else {
		form.Set("AppId", appId)
	}
	form.Set("RoomId", roomId)
	if userId != "" {
		form.Set("UserId", userId)
	}
	if forbiddenInterval != "" {
		form.Set("ForbiddenInterval", forbiddenInterval)
	}
	res, status, err := v.Client.Post("BanRoomUser", nil, form)
	slog.Info("CloseRoom", "res", res, "status", status, "err", err)
	return err
}
