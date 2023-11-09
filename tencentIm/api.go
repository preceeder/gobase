//   File Name:  api.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/7 17:44
//    Change Activity:

package tencentIm

import (
	"errors"
	"fmt"
	"github.com/preceeder/gobase/utils"
	"golang.org/x/exp/rand"
	"log/slog"
	"time"
)

/*
* 发送单聊消息
 *params fromId string 发送者id
 *params toId string 接收者id
 *params content MsgContent 消息内容
 *params cloudCustomData 自定义消息
 *params sendMsgControl []string 消息发送控制选项, NoUnread 不计入未读数、NoLastMsg 不更新绘画列表、 WithMuteNotifications 该条消息的接收方对发送方设置的免打扰选项生效
 *params forbidCallbackControl []string 消息回调禁止开关，只对本条消息有效, ForbidBeforeSendMsgCallback 禁止发消息前回调, ForbidAfterSendMsgCallback 禁止发消息后回调
 *params syncOtherMachine int   1：把消息同步到 From_Account 在线终端和漫游上; 2：消息不同步至 From_Account; 若不填写默认情况下会将消息存 From_Account 漫游
 *params offLineData OfflinePushInfo 离线消息
*/
func SendImMessage(ctx utils.Context, fromId string, toId string, content MsgContent, cloudCustomData string,
	sendMsgControl []string, forbidCallbackControl []string, syncOtherMachine int,
	offLineData *OfflinePushInfo, res BaseResponse) error {

	if len(fromId) == 0 || len(toId) == 0 {
		return errors.New("缺失发送者或接受者")
	}
	// 随机字符串
	var randInt = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

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
	err := SendImRequest(ctx, "SendMsg", message, res)
	if err != nil {
		return err
	} else {
		slog.Info("发送消息response", "response", res.GetResponse())
	}
	return nil
}

/*
* 批量发送单聊消息
 *params fromId string 发送者id
 *params toIds []string 接收者id
 *params content MsgContent 消息内容
 *params cloudCustomData 自定义消息
 *params sendMsgControl []string 消息发送控制选项, NoUnread 不计入未读数、NoLastMsg 不更新绘画列表、 WithMuteNotifications 该条消息的接收方对发送方设置的免打扰选项生效
 *params syncOtherMachine int  0: 根据from_id判断，1：同步，2：不同步
 *params offLineData OfflinePushInfo 离线消息
*/
func SendBatchImMessage(ctx utils.Context, fromId string, toIds []string, content MsgContent, cloudCustomData string,
	sendMsgControl []string, syncOtherMachine int, offLineData *OfflinePushInfo, res BaseResponse) error {

	if len(fromId) == 0 || len(toIds) == 0 {
		return errors.New("缺失发送者或接受者")
	}
	// 随机字符串
	var randInt = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

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
	err := SendImRequest(ctx, "BatchSendMsg", message, res)
	if err != nil {
		return err
	} else {
		slog.Info("发送消息response", "response", res.GetResponse())
	}
	return nil
}

/**
 *消息撤回
 *params fromAccount string 发送者id
 *params toAccount string 接受者id
 *params msgKey string  消息的key
 */
func MsgWithdraw(ctx utils.Context, fromAccount, toAccount, msgKey string) error {
	requestData := map[string]string{
		"From_Account": fromAccount,
		"To_Account":   toAccount,
		"MsgKey":       msgKey,
	}
	res := CommonResponse{}

	err := SendImRequest(ctx, "MsgWithdraw", requestData, &res)
	if err != nil {
		return err
	} else {
		slog.Info("发送消息response", "response", res.GetResponse())
	}
	return nil
}

/**
* 删除会话
 *params fromAccount string 请求删除该 UserID 的会话
 *params toAccount string C2C 会话才赋值，C2C 会话方的 UserID
 *params ToGroupid string G2C 会话才赋值，G2C 会话的群 ID
 *params htype int  会话类型 1：表示 C2C 会话, 2：表示 G2C 会话
 *params ClearRamble int 是否清理漫游消息：1：表示清理漫游消息, 0：表示不清理漫游消息
*/
func DeleteRecentContact(ctx utils.Context, fromAccount, toAccount, toGroupid string, htype, clearRamble int) error {
	requestData := struct {
		FromAccount string `json:"From_Account"`
		ToAccount   string `json:"To_Account,omitempty"`
		ToGroupid   string `json:"ToGroupid,omitempty"`
		Type        int    `json:"Type"`
		ClearRamble int    `json:"ClearRamble,omitempty"`
	}{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		ToGroupid:   toGroupid,
		Type:        htype,
		ClearRamble: clearRamble,
	}

	res := struct {
		CommonResponse
		ErrorDisplay string `json:"ErrorDisplay" mapstructure:"ErrorDisplay"`
	}{}

	err := SendImRequest(ctx, "DeleteRecentContact", requestData, &res)
	if err != nil {
		return err
	} else {
		slog.Info("发送消息response", "response", res.GetResponse())
	}
	return nil
}

/*
*
*拉取历史信息
*正常情况下，分别以会话双方的角度查询消息，结果是一样的。但以下四种情况会导致结果不一样（即会话里的某些消息，其中一方能查询到，另一方查询不到）：
- 会话的其中一方清空了会话的消息记录，即调用了终端的 clearC2CHistoryMessage() 接口。
- 会话的其中一方删除了会话，即调用了终端的 deleteConversation() 接口，或者 Web /小程序/ uni-app 的 deleteConversation 接口，或者服务端的 删除单个会话 的接口且指定了 ClearRamble 的值为1。
- 会话的其中一方删除了部分消息，即调用了终端的 deleteMessages() 接口，或者 Web /小程序/ uni-app 的 deleteMessage 接口。
- 通过 单发单聊消息 或 批量发单聊消息 接口发送的消息，指定了 SyncOtherMachine 值为2，即指定消息不同步到发送方的消息记录
*params operatorAccount string 会话其中一方的 UserID，以该 UserID 的角度去查询消息。同一个会话，分别以会话双方的角度去查询消息，结果可能会不一样，请参考本接口的接口说明
*params peerAccount string 会话的另一方 UserID
*params lastMsgKey string 上一次拉取到的最后一条消息的 MsgKey，续拉时需要填该字段
*params maxCnt int 请求的消息条数
*params minTime int 请求的消息时间范围的最小值（单位：秒）
*params maxTime int 请求的消息时间范围的最大值（单位：秒）
*/
func QueryHistoryMessage(ctx utils.Context, operatorAccount, peerAccount, lastMsgKey string,
	maxCnt int, minTime, maxTime int64) (any, error) {
	requestData := struct {
		OperatorAccount string `json:"Operator_Account"`
		PeerAccount     string `json:"Peer_Account"`
		LastMsgKey      string `json:"LastMsgKey,omitempty"`
		MaxCnt          int    `json:"MaxCnt"`
		MinTime         int64  `json:"MinTime"`
		MaxTime         int64  `json:"MaxTime"`
	}{
		OperatorAccount: operatorAccount,
		PeerAccount:     peerAccount,
		LastMsgKey:      lastMsgKey,
		MaxCnt:          maxCnt,
		MinTime:         minTime,
		MaxTime:         maxTime,
	}

	res := map[string]any{}
	fmt.Println(requestData)
	err := SendImRequest(ctx, "QueryMsg", requestData, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
