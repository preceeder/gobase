/*
File Name:  push.go
Description:
Author:      Chenghu
Date:       2023/9/20 11:58
Change Activity:
*/
package push

import (
	number "github.com/alibabacloud-go/darabonba-number/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	push20160801 "github.com/alibabacloud-go/push-20160801/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/bytedance/sonic"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"strings"
)

var AliPushClient *push20160801.Client
var AliPushConfig AliPush

func InitWithViper(config viper.Viper) {
	//aliConfig := readAliPushConfig(config)
	utils.ReadViperConfig(config, "ali_push", &AliPushConfig)
	_, err := CreateClient(&(AliPushConfig.KeyId), &(AliPushConfig.Secret),
		&(AliPushConfig.EndPoint), &(AliPushConfig.RegionId))
	if err != nil {
		slog.Error("阿里云push创建失败", "error", err.Error())
		panic("阿里云push创建失败：" + err.Error())
	}
}

type AliPush struct {
	KeyId         string `json:"keyId"`
	Secret        string `json:"secret"`
	EndPoint      string `json:"endPoint"`
	RegionId      string `json:"regionId"`
	AppKeyAndroid string `json:"appKeyAndroid"`
	AppKeyIos     string `json:"appKeyIos"`
	Env           string `json:"env"`
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string, endpoint *string, regionId *string) (_result *push20160801.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
		RegionId:        regionId,
		Endpoint:        endpoint,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Push
	//config.Endpoint = endpoint // tea.String("cloudpush.aliyuncs.com")
	_result = &push20160801.Client{}
	_result, _err = push20160801.NewClient(config)
	AliPushClient = _result
	return _result, _err
}

/**
 *  @param userIds []string 用户id 列表
 * @param title string 转为通知的标题
 * @param message map[string]any 发送给用户的消息， 自定义信息
 * @param StoreOffline bool 是否离线推送
 * @param alter bool 是否离线弹窗
 * @param content string 离线弹窗时的内容
 * @param env string PRODUCT | DEV
 */
func GetMessageFormat(ctx utils.Context, userIds []string, title string, message map[string]any, StoreOffline bool,
	alter bool, content string, env string) *push20160801.MassPushRequestPushTask {

	//message := map[string]any{
	//	"type":    "",
	//	"data":    message,
	//}
	extParameters, err := sonic.MarshalString(map[string]any{"push": message})
	if err != nil {
		slog.Error("message json marshal error", "error", err.Error(), "requestId", ctx.RequestId)
		return nil
	}
	body, err := sonic.MarshalString(message)
	if err != nil {
		slog.Error("message json marshal error", "error", err.Error(), "requestId", ctx.RequestId)
		return nil
	}

	pushTask := &push20160801.MassPushRequestPushTask{
		PushType:                       tea.String("MESSAGE"),
		DeviceType:                     tea.String("ALL"),
		StoreOffline:                   tea.Bool(StoreOffline),
		Target:                         tea.String("ACCOUNT"),
		TargetValue:                    tea.String(strings.Join(userIds, ",")),
		Title:                          tea.String(title),
		AndroidNotifyType:              tea.String("VIBRATE"),
		AndroidOpenType:                tea.String("APPLICATION"),
		AndroidActivity:                tea.String(""),
		AndroidNotificationBarType:     tea.Int32(50),
		AndroidNotificationBarPriority: tea.Int32(0),
		AndroidExtParameters:           tea.String(extParameters),
		AndroidNotificationChannel:     tea.String("静默提醒"),
		IOSApnsEnv:                     tea.String(env),
		IOSSilentNotification:          tea.Bool(true),
		IOSMutableContent:              tea.Bool(true),
		IOSExtParameters:               tea.String(extParameters),
		IOSBadgeAutoIncrement:          tea.Bool(true),
		Body:                           tea.String(body),
	}
	if alter {
		pushTask.AndroidRemind = tea.Bool(true)
		pushTask.AndroidPopupActivity = tea.String("")
		pushTask.AndroidPopupTitle = tea.String(title)
		pushTask.AndroidPopupBody = tea.String(content)
		pushTask.IOSRemind = tea.Bool(true)
		pushTask.IOSRemindBody = tea.String(content)
	} else {
		pushTask.AndroidRemind = tea.Bool(false)
	}
	return pushTask
}

/*
 * @param pushTask *push20160801.MassPushRequestPushTask 先调用GetMessageFormat 拿到结果就是这里的参数
 * @param appKey string   由于android 和ios 可能不一样所以这里需要给个参数
 */
func PushMessage(ctx utils.Context, pushTask *push20160801.MassPushRequestPushTask, appKey string) {
	request := &push20160801.MassPushRequest{
		AppKey:   number.ParseLong(&appKey),
		PushTask: []*push20160801.MassPushRequestPushTask{pushTask},
	}

	// request.pushTask = new Push20160801.MassPushRequest.pushTask{};
	_, _err := AliPushClient.MassPush(request)
	if _err != nil {
		slog.Error("阿里云 推送消息失败", "error", _err.Error(), "requestId", ctx.RequestId)
		return
	}
}
