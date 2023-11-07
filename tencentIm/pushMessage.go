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

/*
* 发送单聊消息
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
		res = &CommonResponse{}
	}
	//fmt.Println(sonic.ConfigFastest.MarshalToString(message))
	err := SendImRequest(ctx, serverName, message, res)
	if err != nil {
		return err
	} else {
		slog.Info("发送消息response", "response", res.GetResponse())
	}
	return nil
}

/*
* 批量发送单聊消息
params serverName string url的映射 ApiMap
params fromId string 发送者id
params toIds []string 接收者id
params content MsgContent 消息内容
params cloudCustomData 自定义消息
params sendMsgControl []string 消息发送控制选项, NoUnread 不计入未读数、NoLastMsg 不更新绘画列表、 WithMuteNotifications 该条消息的接收方对发送方设置的免打扰选项生效
params syncOtherMachine int  0: 根据from_id判断，1：同步，2：不同步
params offLineData OfflinePushInfo 离线消息
*/
func SendBatchImMessage(ctx utils.Context, serverName string, fromId string, toIds []string, content MsgContent, cloudCustomData string,
	sendMsgControl []string, syncOtherMachine int, offLineData *OfflinePushInfo, res BaseResponse) error {

	if len(fromId) == 0 || len(toIds) == 0 {
		return errors.New("缺失发送者或接受者")
	}
	// 随机字符串
	var randInt = rand.New(rand.NewSource(time.Now().UnixNano()))

	var message = BatchMessage{
		SyncOtherMachine: syncOtherMachine,
		MsgLifeTime:      3600 * 24 * 7,
		FromAccount:      fromId,
		ToAccount:        toIds,
		MsgRandom:        randInt.Intn(1000000),
		SendMsgControl:   sendMsgControl,
		CloudCustomData:  cloudCustomData,
		OfflinePushInfo:  offLineData,
		MsgBody:          []MsgBody{{MsgType: content.GetMsgType(), MsgContent: content}},
	}
	if res == nil {
		res = &BatchCommonResponse{}
	}
	err := SendImRequest(ctx, serverName, message, res)
	if err != nil {
		return err
	} else {
		slog.Info("发送消息response", "response", res.GetResponse())
	}
	return nil
}
