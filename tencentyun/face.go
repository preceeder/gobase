package tencentyun

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/preceeder/gobase/db/redisDb"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"
)

var TencentFaceClient *TencentFace
var TencentFaceConfig tencentFaceConfig

type tencentFaceConfig struct {
	AppId        string `json:"appId"`
	ServerSecret string `json:"serverSecret"`
}

func InitWithViper(config viper.Viper) {
	utils.ReadViperConfig(config, "tencentyun", &TencentFaceConfig)
	InitTencentFace(TencentFaceConfig.AppId, TencentFaceConfig.ServerSecret)
	InitTencentFaceId(TencentFaceConfig.AppId, TencentFaceConfig.ServerSecret)
}

func InitTencentFace(appid string, serverSecret string) {
	TencentFaceClient = NewTencentFaceClient(appid, serverSecret)
}

func NewTencentFaceClient(appid string, serverSecret string) *TencentFace {
	return &TencentFace{
		AppId:        appid,
		ServerSecret: serverSecret,
		RestyClient:  resty.New().SetTimeout(time.Duration(5 * time.Second)),
	}
}

type TencentFace struct {
	AppId        string
	ServerSecret string
	RestyClient  *resty.Client
}

func (tf TencentFace) getTencentFaceAccessToken(ctx utils.Context) (accessToken string) {
	// 获取access_token
	url := "https://miniprogram-kyc.tencentcloudapi.com/api/oauth2/access_token"
	params := map[string]string{
		"app_id":     tf.AppId,
		"secret":     tf.ServerSecret,
		"grant_type": "client_credential",
		"version":    "1.0.0",
	}
	resp, err := tf.RestyClient.R().SetQueryParams(params).Get(url)
	if err != nil {
		slog.Error("访问腾讯云access_token接口异常", "err", err, "requestId", ctx.RequestId)
		return
	}

	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		slog.Error("解析腾讯云access_token接口响应体失败", "err", err, "requestId", ctx.RequestId)
		return
	}
	if result["code"] != "0" {
		slog.Error("访问腾讯云access_token接口失败", "err", err, "requestId", ctx.RequestId)
		return
	}
	str, ok := result["access_token"]
	if !ok {
		slog.Error("腾讯云access_token接口解析access_token参数失败", "err", err, "requestId", ctx.RequestId)
	}
	accessToken = str.(string)
	return
}

func (tf TencentFace) getTencentFaceSignTicket(ctx utils.Context) (signTicket string) {
	// 获取签名钥匙
	cmd, err := redisDb.Do(ctx, map[string]any{"cmd": "get {{redisKey}}"}, map[string]any{"redisKey": "tencentFaceSignTicket"})
	if err != nil {
		slog.Error("连接redis失败", "err", err, "requestId", ctx.RequestId)
	} else {
		signTicket, err = cmd.Text()
		if err == nil {
			return
		}
	}
	accessToken := tf.getTencentFaceAccessToken(ctx)
	if accessToken == "" {
		return
	}
	url := "https://miniprogram-kyc.tencentcloudapi.com/api/oauth2/api_ticket"
	params := map[string]string{
		"app_id":       tf.AppId,
		"access_token": accessToken,
		"type":         "SIGN",
		"version":      "1.0.0",
	}
	resp, err := tf.RestyClient.R().SetQueryParams(params).Get(url)
	if err != nil {
		slog.Error("访问腾讯云SignTicket接口异常", "err", err, "requestId", ctx.RequestId)
		return
	}
	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		slog.Error("解析腾讯云SignTicket接口响应体失败", "err", err, "requestId", ctx.RequestId)
		return
	}
	if result["code"] != "0" {
		slog.Error("访问腾讯云SignTicket接口失败", "err", err, "requestId", ctx.RequestId)
		return
	}
	tickets, ok := result["tickets"]
	if !ok {
		slog.Error("腾讯云SignTicket接口解析tickets参数失败", "err", err, "requestId", ctx.RequestId)
		return
	}
	signTicket, ok = tickets.([]map[string]string)[0]["value"]
	if !ok {
		slog.Error("腾讯云SignTicket接口解析signTicket参数失败", "err", err, "requestId", ctx.RequestId)
	}
	_, err = redisDb.Do(ctx, map[string]any{"cmd": "set {{redisKey}} {{value}}", "key": "{{redisKey}}", "exp": time.Minute * 20}, map[string]any{"redisKey": "tencentFaceSignTicket", "value": signTicket})
	if err != nil {
		slog.Error("连接redis失败", "err", err, "requestId", ctx.RequestId)
	}
	return
}

func (tf TencentFace) cookTencentSign(ctx utils.Context, userId string, nonce string) (sign string) {
	// 制作签名
	signTicket := tf.getTencentFaceSignTicket(ctx)
	if signTicket == "" {
		slog.Error("获取signTicket失败", "requestId", ctx.RequestId)
		return
	}

	type KVPair struct {
		Key string
		Val string
	}
	signParams := map[string]string{
		"WBappid": tf.AppId,
		"userId":  userId,
		"version": "1.0.0",
		"ticket":  signTicket,
		"nonce":   nonce,
	}
	values := maps.Values(signParams)
	slices.Sort(values)
	str := strings.Join(values, "")
	s := sha1.New()
	s.Write([]byte(str))
	sign = hex.EncodeToString(s.Sum(nil))
	return
}

type TencentFaceIdResult struct {
	FaceID        string `json:"faceId"`
	AgreementNo   string `json:"agreementNo"`
	OpenAPINonce  string `json:"openApiNonce"`
	OpenAPIUserID string `json:"openApiUserId"`
	OpenAPISign   string `json:"openApiSign"`
}

func (tf TencentFace) GetTencentFaceId(ctx utils.Context, userId string, imgBase64 string, orderNo string) any {
	// 人脸核身faceId获取
	nonce := utils.GenterWithoutRepetitionStr(32)
	sign := tf.cookTencentSign(ctx, userId, nonce)
	if sign == "" {
		return nil
	}
	data := map[string]string{
		"webankAppId":     tf.AppId,
		"orderNo":         orderNo,
		"userId":          userId,
		"sourcePhotoStr":  imgBase64,
		"sourcePhotoType": "2",
		"version":         "1.0.0",
		"sign":            sign,
		"nonce":           nonce,
	}
	resp, err := tf.RestyClient.R().SetHeader("Content-Type", "application/json").SetFormData(data).Post(fmt.Sprintf("https://miniprogram-kyc.tencentcloudapi.com/api/server/getfaceid?orderNo=%s", orderNo))
	if err != nil {
		slog.Error("获取人脸认证token失败", "err", err, "requestId", ctx.RequestId)
		return nil
	}
	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	if result["code"] != "0" {
		slog.Error("访问腾讯云faceId接口失败", "err", err, "requestId", ctx.RequestId)
		return nil
	}
	results, ok := result["result"]
	if !ok {
		slog.Error("腾讯云faceId接口解析faceId参数失败", "err", err, "requestId", ctx.RequestId)
		return nil
	}
	faceId, ok := results.(map[string]string)["faceId"]
	if !ok {
		slog.Error("腾讯云faceId接口解析faceId参数失败", "err", err, "requestId", ctx.RequestId)
		return nil
	}
	return TencentFaceIdResult{
		FaceID:        faceId,
		AgreementNo:   orderNo,
		OpenAPINonce:  nonce,
		OpenAPIUserID: userId,
		OpenAPISign:   sign,
	}
}

func (tf TencentFace) GetTencentFaceResult(ctx utils.Context, userId string, orderNo string) (similarity float64, ok bool) {
	// 人脸核身结果查询
	nonce := utils.GenterWithoutRepetitionStr(32)
	sign := tf.cookTencentSign(ctx, userId, nonce)
	if sign == "" {
		return
	}
	data := map[string]string{
		"appId":   tf.AppId,
		"version": "1.0.0",
		"nonce":   nonce,
		"orderNo": orderNo,
		"sign":    sign,
	}
	url := fmt.Sprintf("https://miniprogram-kyc.tencentcloudapi.com/api/v2/base/queryfacerecord?orderNo=%s", orderNo)
	resp, err := tf.RestyClient.R().SetFormData(data).Post(url)
	if err != nil {
		return
	}
	var result map[string]any
	err = json.Unmarshal(resp.Body(), &result)
	if result["code"] != "0" {
		slog.Error("访问腾讯云人脸核身结果查询接口失败", "err", err, "requestId", ctx.RequestId)
		return
	}
	results, ok := result["result"]
	if !ok {
		slog.Error("访问腾讯云人脸核身结果查询接口获取similarity失败", "err", err, "requestId", ctx.RequestId)
		return
	}
	res, ok := results.(map[string]string)["similarity"]
	if !ok {
		slog.Error("访问腾讯云人脸核身结果查询接口获取similarity失败", "err", err, "requestId", ctx.RequestId)
		return
	}
	similarity, err = strconv.ParseFloat(res, 64)
	if err != nil {
		return
	}
	return similarity, true
}
