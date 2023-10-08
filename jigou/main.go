/*
File Name:  main.py
Description:
Author:      Chenghu
Date:       2023/8/22 09:36
Change Activity:
*/
package jigou

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/go-resty/resty/v2"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"github.com/zegoim/zego_server_assistant/token/go/src/token04"
	"log/slog"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// AppId = ""
// ServerSecret = ""

var JiGouClient *JiGou
var JiGouConfig jigouConfig

const (
	JiGouPublicUrl             string = "https://rtc-api.zego.im"
	JiGouCloseRoom             string = "https://rtc-api.zego.im/?Action=CloseRoom"         // 关闭房间
	JiGouDescribeUsers         string = "https://rtc-api.zego.im/?Action=DescribeUsers"     // 查询用户状态
	JiGouSendCustomCommand     string = "https://rtc-api.zego.im/?Action=SendCustomCommand" // 推送自定义消息
	JiGouDescribeUserNum       string = "https://rtc-api.zego.im/?Action=DescribeUserNum"
	JiGouGenerateIdentifyToken string = "https://rtc-api.zego.im/?Action=GenerateIdentifyToken"
)

type jigouConfig struct {
	AppId          string `json:"appId"`
	ServerSecret   string `json:"serverSecret"`
	CallBackSecret string `json:"callbackSecret"`
}

func readJigouConfig(v viper.Viper) (jg jigouConfig) {
	jigou := v.Sub("jigou")
	if jigou == nil {
		fmt.Printf("jigou config is nil")
		os.Exit(1)
	}
	jg = jigouConfig{}
	err := jigou.Unmarshal(&jg)
	if err != nil {
		fmt.Printf("jigou config read error: " + err.Error())
		os.Exit(1)
	}
	return
}

func InitWithViper(config viper.Viper) {
	jigouConfig := readJigouConfig(config)
	InitJiGou(jigouConfig.AppId, jigouConfig.ServerSecret)
}

func InitJiGou(appid string, serverSecret string) {
	JiGouClient = &JiGou{
		AppId:        appid,
		ServerSecret: serverSecret,
		RestyClient:  resty.New().SetTimeout(time.Duration(5 * time.Second)),
	}
}

func NewJiGouClient(appid string, serverSecret string) *JiGou {
	return &JiGou{
		AppId:        appid,
		ServerSecret: serverSecret,
		RestyClient:  resty.New().SetTimeout(time.Duration(5 * time.Second)),
	}
}

type JiGou struct {
	AppId        string
	ServerSecret string
	RestyClient  *resty.Client
}

// 生成签名
// Signature=md5(AppId + SignatureNonce + ServerSecret + Timestamp)
func generateSignature(appId string, serverSecret, signatureNonce string, timeStamp int64) string {
	data := fmt.Sprintf("%s%s%s%d", appId, signatureNonce, serverSecret, timeStamp)

	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (j JiGou) GetPublicParams() map[string]string {
	publicParams := map[string]string{}
	timestamp := time.Now().Unix()
	// 生成16进制随机字符串(16位)
	nonce := make([]byte, 8)
	_, _ = rand.Read(nonce)
	hexNonce := hex.EncodeToString(nonce)
	// 生成签名
	signature := generateSignature(j.AppId, j.ServerSecret, hexNonce, timestamp)
	publicParams["AppId"] = j.AppId
	//公共参数中的随机数和生成签名的随机数要一致
	publicParams["SignatureNonce"] = hexNonce
	publicParams["SignatureVersion"] = "2.0"
	//公共参数中的时间戳和生成签名的时间戳要一致
	publicParams["Timestamp"] = fmt.Sprintf("%d", timestamp)
	publicParams["Signature"] = signature
	return publicParams
}

func (j JiGou) Get(ctx utils.Context, url string, params url.Values, resBody any) error {
	publicParams := j.GetPublicParams()
	//finalParams := maputil.Merge[string, string](params, publicParams)
	_, err := j.RestyClient.R().EnableTrace().SetHeader("Accept", "application/json").
		SetResult(resBody).SetQueryParams(publicParams).SetQueryParamsFromValues(params).Get(url)
	if err != nil {
		slog.Error("jigou request", "error", err.Error(), "requestId", ctx.RequestId)
		return err
	}
	//var tempResp = &map[string]any{}
	//_ = sonic.Unmarshal(resp.Body(), tempResp)
	return nil
}

// 关闭房间

func (j JiGou) CloseRoom(ctx utils.Context, RoomId string) (PublicResponse, error) {
	//roomId := map[string]string{
	//	"RoomId[]": RoomId,
	//}
	roomId := url.Values{
		"RoomId": []string{RoomId},
	}

	resBody := &PublicResponse{}
	err := j.Get(ctx, JiGouCloseRoom, roomId, resBody)
	return *resBody, err
}

// 获取房间的人数

func (j JiGou) GetRoomNumbers(ctx utils.Context, roomId string) (RoomNumbers, error) {
	params := url.Values{
		"RoomId[]": []string{roomId},
	}
	resBody := &RoomNumbers{}
	err := j.Get(ctx, JiGouDescribeUserNum, params, resBody)
	//fmt.Printf("%#v\n", resBody)
	return *resBody, err
}

// 发送自定义消息

func (j JiGou) SendCustomCommand(ctx utils.Context, roomId string, fromUserId string, toUserId []string, message string) (SendCustomCommand, error) {
	//params := map[string]string{
	//	"RoomId":         roomId,
	//	"FromUserId":     fromUserId,
	//	"ToUserId[]":     toUserId[0],
	//	"MessageContent": message,
	//}
	params := url.Values{
		"RoomId":         []string{roomId},
		"FromUserId":     []string{fromUserId},
		"ToUserId[]":     toUserId,
		"MessageContent": []string{message},
	}
	resBody := &SendCustomCommand{}
	err := j.Get(ctx, JiGouDescribeUserNum, params, resBody)
	//fmt.Printf("%#v\n", resBody)
	return *resBody, err
}

// 获取 音视频流审核鉴权 Token

func (j JiGou) GenerateIdentifyToken(ctx utils.Context) (GenerateIdentifyToken, error) {
	resBody := &GenerateIdentifyToken{}
	err := j.Get(ctx, JiGouGenerateIdentifyToken, nil, resBody)
	return *resBody, err
}

// 获取jigou权限认证token   客户端使用
func (j JiGou) GetToken(ctx utils.Context, userId, roomId string) (string, error) {
	var effectiveTimeInSeconds int64 = 3600 // token 的有效时长，单位：秒
	//业务权限认证配置，可以配置多个权限位
	privilege := make(map[int]int)
	privilege[token04.PrivilegeKeyLogin] = token04.PrivilegeEnable   // 有房间登录权限
	privilege[token04.PrivilegeKeyPublish] = token04.PrivilegeEnable // 无推流权限
	//token业务扩展配置
	payloadData := &RtcRoomPayLoad{
		RoomId:       roomId,
		Privilege:    privilege,
		StreamIdList: nil,
	}
	payload, err := sonic.MarshalString(payloadData)
	if err != nil {
		slog.Error("GetToken error", "error", err.Error(), "requestId", ctx.RequestId)
		return "", err
	}
	//生成token
	appid, _ := strconv.Atoi(j.AppId)
	token, err := token04.GenerateToken04(uint32(appid), userId, j.ServerSecret, effectiveTimeInSeconds, payload)
	if err != nil {
		fmt.Println(err)
	}
	return token, nil
}

// 回调参数校验
func CallDataCheck(timestamp, nonce, signature string) bool {
	data := []string{JiGouConfig.CallBackSecret, timestamp, nonce}
	slices.Sort(data)
	chd := cryptor.Sha1(strings.Join(data, ""))
	if chd == signature {
		return true
	}
	return false
}
