//   File Name:  pushMessage.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/6 13:47
//    Change Activity:

package tencentIm

import (
	"errors"
	"github.com/preceeder/gobase/utils"
	"log/slog"
	"math/rand"
	"time"
)

type Message struct {
	//1：把消息同步到 From_Account 在线终端和漫游上
	//2：消息不同步至 From_Account
	//若不填写默认情况下会将消息存 From_Account 漫游
	SyncOtherMachine int `json:"SyncOtherMachine,omitempty"`
	//消息发送方 UserID（用于指定发送消息方账号）
	FromAccount string `json:"From_Account"`
	//消息接收方 UserID
	ToAccount string `json:"To_Account"`
	//消息离线保存时长（单位：秒），最长为7天（604800秒）
	//若设置该字段为0，则消息只发在线用户，不保存离线
	//若设置该字段超过7天（604800秒），仍只保存7天
	//若不设置该字段，则默认保存7天
	MsgLifeTime int `json:"MsgLifeTime,omitempty"`
	//消息序列号（32位无符号整数），后台会根据该字段去重及进行同秒内消息的排序，详细规则请看本接口的功能说明。若不填该字段，则由后台填入随机数
	MsgSeq int `json:"MsgSeq,omitempty"`
	//消息随机数（32位无符号整数），后台用于同一秒内的消息去重。请确保该字段填的是随机
	MsgRandom int `json:"MsgRandom"`
	//消息回调禁止开关，只对本条消息有效，ForbidBeforeSendMsgCallback 表示禁止发消息前回调，ForbidAfterSendMsgCallback 表示禁止发消息后回调
	ForbidCallbackControl []string `json:"ForbidCallbackControl,omitempty"`
	//消息发送控制选项，是一个 String 数组，只对本条消息有效。"NoUnread"表示该条消息不计入未读数。
	//"NoLastMsg"表示该条消息不更新会话列表。
	//"WithMuteNotifications"表示该条消息的接收方对发送方设置的免打扰选项生效（默认不生效）。
	//"NoMsgCheck"表示开启云端审核后，该条消息不送审。
	//示例：
	//"SendMsgControl": ["NoUnread","NoLastMsg","WithMuteNotifications","NoMsgCheck"]
	SendMsgControl []string `json:"SendMsgControl,omitempty"`
	//消息内容，具体格式请参考 消息格式描述（注意，一条消息可包括多种消息元素，MsgBody 为 Array 类型）
	MsgBody []MsgBody `json:"MsgBody"`

	//消息自定义数据（云端保存，会发送到对端，程序卸载重装后还能拉取到）
	CloudCustomData string `json:"CloudCustomData,omitempty"`
	//该条消息是否支持消息扩展，0为不支持，1为支持。
	SupportMessageExtension int `json:"SupportMessageExtension,omitempty"`
	// 离线推送信息配置，具体可参考 消息格式描述
	OfflinePushInfo *OfflinePushInfo `json:"OfflinePushInfo,omitempty"`
	//该条消息是否需要已读回执，0为不需要，1为需要，默认为0
	IsNeedReadReceipt int `json:"IsNeedReadReceipt,omitempty"`
}

type MsgBody struct {
	//TIM 消息对象类型，目前支持的消息对象包括：
	//TIMTextElem（文本消息）
	//TIMLocationElem（位置消息）
	//TIMFaceElem（表情消息）
	//TIMCustomElem（自定义消息）
	//TIMSoundElem（语音消息）
	//TIMImageElem（图像消息）
	//TIMFileElem（文件消息）
	//TIMVideoFileElem（视频消息）
	MsgType string `json:"MsgType"`
	//对于每种 MsgType 用不同的 MsgContent 格式，具体可参考 消息格式描述
	MsgContent MsgContent `json:"MsgContent"`
}

type OfflinePushInfo struct {
	//0：表示推送
	//1：表示不离线推送
	PushFlag int `json:"PushFlag"`
	// 离线推送标题。该字段为 iOS 和 Android 共用。
	Title string `json:"Title"`
	//离线推送内容。该字段会覆盖上面各种消息元素 TIMMsgElement 的离线推送展示文本。
	//若发送的消息只有一个 TIMCustomElem 自定义消息元素，
	//该 Desc 字段会覆盖 TIMCustomElem 中的 Desc 字段。
	//如果两个 Desc 字段都不填，将收不到该自定义消息的离线推送。
	Desc string `json:"Desc"`
	//离线推送透传内容。由于国内各 Android 手机厂商的推送平台要求各不一样，请保证此字段为 JSON 格式，否则可能会导致收不到某些厂商的离线推送。
	Ext         string      `json:"Ext"`
	AndroidInfo AndroidInfo `json:"AndroidInfo"`
	ApnsInfo    ApnsInfo    `json:"ApnsInfo"`
}
type AndroidInfo struct {
	//Android 离线推送声音文件路径。
	Sound string `json:"Sound,omitempty"`
	//Android通知栏样式，“0”代表默认样式，“1”代表大文本样式，不填默认为0。仅对
	//华为/荣耀/OPPO生效。
	PushStyle int `json:"PushStyle,omitempty"`
	//华为手机 EMUI 10.0 及以上的通知渠道字段。
	//该字段不为空时，会覆盖控制台配置的 ChannelID 值；
	//该字段为空时，不会覆盖控制台配置的 ChannelID 值。
	HuaWeiChannelID string `json:"HuaWeiChannelID,omitempty"`
	//华为推送通知消息分类，取值为 LOW、NORMAL，不填默认为 NORMAL。
	HuaWeiImportance string `json:"HuaWeiImportance,omitempty"`
	//在控制台配置华为推送为“打开应用内指定页面”的前提下，
	//传“1”表示将透传内容 Ext 作为 Intent 的参数,
	//“0”表示将透传内容 Ext 作为 Action 参数。
	//不填默认为0。两种传参区别可参见 华为推送文档。
	ExtAsHuaweiIntentParam int `json:"ExtAsHuaweiIntentParam,omitempty"`
	//华为手机用来标识消息类型，该字段不为空时，会覆盖控制台配置的 category 值；
	//该字段为空时，不会覆盖控制台配置的 category 值。详见 category 描述
	HuaWeiCategory string `json:"HuaWeiCategory,omitempty"`
	//华为推送通知栏消息右侧小图标URL，URL必须使用HTTPS协议，取值样例：https://example.com/image.png。
	//图片文件须小于512KB，规格建议为40dp x 40dp，弧角大小为8dp。
	//超出建议规格的图片会存在图片压缩或图片显示不全的情况。图片格式建议使用JPG/JPEG/PNG
	HuaWeiImage string `json:"HuaWeiImage,omitempty"`
	//荣耀推送通知栏消息右侧大图标URL，URL必须使用HTTPS协议， 取值样例：https://example.com/image.png。
	//图标文件须小于512KB，图标建议规格大小：40dp x 40dp，
	//弧角大小为8dp，超出建议规格大小的图标会存在图片压缩或显示不全的情况。
	HonorImage string `json:"HonorImage,omitempty"`
	//荣耀推送通知消息分类，取值为 LOW、NORMAL，不填默认为 NORMAL。
	HonorImportance string `json:"HonorImportance,omitempty"`
	//Google 推送通知栏消息右侧图标URL，图片资源不超过1M，支持JPG/JPEG/PNG格式，取值样例：https://example.com/image.png
	GoogleImage string `json:"GoogleImage"`
	//小米手机 MIUI 10 及以上的通知类别（Channel）适配字段。
	//该字段不为空时，会覆盖控制台配置的 ChannelID 值；
	//该字段为空时，不会覆盖控制台配置的 ChannelID 值。
	XiaoMiChannelID string `json:"XiaoMiChannelID,omitempty"`
	//OPPO 手机 Android 8.0 及以上的 NotificationChannel 通知适配字段。
	//该字段不为空时，会覆盖控制台配置的 ChannelID 值；
	//该字段为空时，不会覆盖控制台配置的 ChannelID 值。
	OPPOChannelID string `json:"OPPOChannelID,omitempty"`
	//Google 手机 Android 8.0 及以上的通知渠道字段。
	//Google 推送新接口（上传证书文件）支持 channel id，旧接口（填写服务器密钥）不支持。
	GoogleChannelID string `json:"GoogleChannelID,omitempty"`
	//VIVO 手机用来标识消息类型，该字段不为空时，会覆盖控制台配置的 category 值；
	//该字段为空时，不会覆盖控制台配置的 category 值。
	//详见 category 描述
	VIVOClassification int `json:"VIVOClassification,omitempty"`
}
type ApnsInfo struct {
	Sound string `json:"Sound,omitempty"`
	//这个字段缺省或者为0表示需要计数，为1表示本条消息不需要计数，即右上角图标数字不增加。
	BadgeMode int `json:"BadgeMode,omitempty"`
	//该字段用于标识 APNs 推送的标题，若填写则会覆盖最上层 Title。
	Title string `json:"Title,omitempty"`
	//该字段用于标识 APNs 推送的子标题。
	SubTitle string `json:"SubTitle,omitempty"`
	//该字段用于标识 APNs 携带的图片地址，当客户端拿到该字段时，可以通过下载图片资源的方式将图片展示在弹窗上。
	Image string `json:"Image,omitempty"`
	//为1表示开启 iOS 10+ 的推送扩展，默认为0。
	MutableContent int `json:"MutableContent,omitempty"`
}

/*
*
params serverName string url的映射 ApiMap
params fromId string 发送者id
params toId string 接收者id
params content MsgContent 消息内容
params cloudCustomData 自定义消息
params sendMsgControl []string 消息发送控制选项, NoUnread 不计入未读数、NoLastMsg 不更新绘画列表、 WithMuteNotifications 该条消息的接收方对发送方设置的免打扰选项生效
params forbidCallbackControl []string 消息回调禁止开关，只对本条消息有效, ForbidBeforeSendMsgCallback 禁止发消息前回调, ForbidAfterSendMsgCallback 禁止发消息后回调
params syncOtherMachine int  0: 根据from_id判断，1：同步，2：不同步
params offLineData OfflinePushInfo 离线消息
*/
func SendImMessage(ctx utils.Context, serverName string, fromId string, toId string, content MsgContent, cloudCustomData string,
	sendMsgControl []string, forbidCallbackControl []string, syncOtherMachine int,
	offLineData *OfflinePushInfo, res BaseResponse) error {

	if len(fromId) == 0 || len(toId) == 0 {
		return errors.New("缺失发送者或接受者")
	}
	// 随机字符串
	var randInt = rand.New(rand.NewSource(time.Now().UnixNano()))

	var message = Message{
		SyncOtherMachine:      syncOtherMachine,
		MsgLifeTime:           3600 * 24 * 7,
		FromAccount:           fromId,
		ToAccount:             toId,
		MsgRandom:             randInt.Intn(1000000),
		ForbidCallbackControl: forbidCallbackControl,
		SendMsgControl:        sendMsgControl,
		CloudCustomData:       cloudCustomData,
		OfflinePushInfo:       offLineData,
		MsgBody:               []MsgBody{{MsgType: content.GetMsgType(), MsgContent: content}},
	}
	if res == nil {
		res = CommonResponse{}
	}
	//fmt.Println(sonic.ConfigFastest.MarshalToString(message))
	err := SendImRequest(ctx, serverName, message, &res)
	if err != nil {
		return err
	} else {
		slog.Info("发送消息response", "response", res.GetResponse())
	}
	return nil
}
