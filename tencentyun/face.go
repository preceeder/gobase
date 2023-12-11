package tencentyun

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/fanjindong/go-cache"
	"github.com/go-resty/resty/v2"
	"github.com/preceeder/gobase/db/dcache"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
	"log/slog"
	"slices"
	"strings"
	"time"
)

var TencentFaceClient *TencentFace
var TencentFaceConfig tencentFaceConfig

type tencentFaceConfig struct {
	Api struct {
		SecretId  string `json:"secretId"`
		SecretKey string `json:"secretKey"`
	} `json:"api"`
	Face struct {
		AppId        string `json:"appId"`
		ServerSecret string `json:"serverSecret"`
	} `json:"face"`
}

func InitWithViper(config viper.Viper) {
	utils.ReadViperConfig(config, "tencentyun", &TencentFaceConfig)
	InitTencentFace(TencentFaceConfig.Face.AppId, TencentFaceConfig.Face.ServerSecret)
	InitTencentFaceId(TencentFaceConfig.Api.SecretId, TencentFaceConfig.Api.SecretKey)
}

func InitTencentFace(appid string, serverSecret string) {
	TencentFaceClient = NewTencentFaceClient(appid, serverSecret)
}

func NewTencentFaceClient(appid string, serverSecret string) *TencentFace {
	return &TencentFace{
		AppId:        appid,
		ServerSecret: serverSecret,
		RestyClient:  resty.New().SetHeader("Content-Type", "application/json").SetTimeout(time.Duration(5 * time.Second)),
	}
}

type TencentFace struct {
	AppId        string
	ServerSecret string
	RestyClient  *resty.Client
}

type FaceAccessTokenResponse struct {
	Code            string `json:"code"`
	Msg             string `json:"msg"`
	TransactionTime string `json:"transactionTime"`
	AccessToken     string `json:"access_token"`
	ExpireTime      string `json:"expire_time"`
	ExpireIn        int    `json:"expire_in"`
}

func (f FaceAccessTokenResponse) String() string {
	res, _ := sonic.ConfigFastest.MarshalToString(f)
	return res
}

func (tf TencentFace) getTencentFaceAccessToken(ctx utils.Context) (accessToken string, err error) {
	// 获取access_token
	url := "https://miniprogram-kyc.tencentcloudapi.com/api/oauth2/access_token"
	params := map[string]string{
		"app_id":     tf.AppId,
		"secret":     tf.ServerSecret,
		"grant_type": "client_credential",
		"version":    "1.0.0",
	}
	var fat = FaceAccessTokenResponse{}
	_, err = tf.RestyClient.R().SetResult(&fat).SetQueryParams(params).Get(url)
	if err != nil {
		slog.Error("访问腾讯云access_token接口异常", "err", err, "requestId", ctx.RequestId)
		return
	}
	if fat.Code != "0" {
		slog.Error("访问腾讯云access_token接口失败", "err", err, "result", fat, "requestId", ctx.RequestId)
		err = errors.New(fat.String())
		return
	}
	accessToken = fat.AccessToken
	return
}

type FaceSignTicketResponse struct {
	Code            string `json:"code"`
	Msg             string `json:"msg"`
	TransactionTime string `json:"transactionTime"`
	Tickets         []struct {
		Value      string `json:"value"`
		ExpireIn   int    `json:"expire_in"`
		ExpireTime string `json:"expire_time"`
	} `json:"tickets"`
}

func (f FaceSignTicketResponse) String() string {
	res, _ := sonic.ConfigFastest.MarshalToString(f)
	return res
}
func (tf TencentFace) getTencentFaceSignTicket(ctx utils.Context) (signTicket string, err error) {
	// 获取签名钥匙
	cacheKey := "BASE_tencentFaceSignTicket"
	sTicket, ok := dcache.GoCache.Get(cacheKey)
	if ok && len(sTicket.(string)) > 0 {
		signTicket = sTicket.(string)
		return
	}

	accessToken, err := tf.getTencentFaceAccessToken(ctx)
	if accessToken == "" || err != nil {
		return
	}
	url := "https://miniprogram-kyc.tencentcloudapi.com/api/oauth2/api_ticket"
	params := map[string]string{
		"app_id":       tf.AppId,
		"access_token": accessToken,
		"type":         "SIGN",
		"version":      "1.0.0",
	}
	var signTi = FaceSignTicketResponse{}
	_, err = tf.RestyClient.R().SetResult(&signTi).SetQueryParams(params).Get(url)
	if err != nil {
		slog.Error("访问腾讯云SignTicket接口异常", "err", err, "requestId", ctx.RequestId)
		return
	}

	if signTi.Code != "0" {
		slog.Error("访问腾讯云SignTicket接口失败", "err", err, "result", signTi, "requestId", ctx.RequestId)
		err = errors.New(signTi.String())
		return
	}

	signTicket = signTi.Tickets[0].Value
	ok = dcache.GoCache.Set(cacheKey, signTicket, cache.WithEx(time.Second*60*20))

	return
}

// 获取face id的时候 需要userId, 不需要orderNo
// 获取结果的时候 需要 orderNo 不需要 userId
func (tf TencentFace) cookTencentSign(ctx utils.Context, userId, orderNo, nonce string) (sign string, err error) {
	// 制作签名
	signTicket, err := tf.getTencentFaceSignTicket(ctx)
	if signTicket == "" || err != nil {
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

	if len(userId) > 0 {
		signParams["userId"] = userId
	} else if len(orderNo) > 0 {
		signParams["orderNo"] = orderNo
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

type TencentFaceIdResponse struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	BizSeqNo string `json:"bizSeqNo"`
	Result   struct {
		BizSeqNo        string `json:"bizSeqNo"`
		TransactionTime string `json:"transactionTime"`
		OrderNo         string `json:"orderNo"`
		FaceID          string `json:"faceId"`
		Success         bool   `json:"success"`
	} `json:"result"`
	TransactionTime string `json:"transactionTime"`
}

func (f TencentFaceIdResponse) String() string {
	res, _ := sonic.ConfigFastest.MarshalToString(f)
	return res
}

func (tf TencentFace) GetTencentFaceId(ctx utils.Context, userId string, imgBase64 string, orderNo string) (any, error) {
	// 人脸核身faceId获取
	nonce := utils.GenterWithoutRepetitionStr(32)
	sign, err := tf.cookTencentSign(ctx, userId, "", nonce)
	if sign == "" {
		return nil, err
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
	var tfir = TencentFaceIdResponse{}
	_, err = tf.RestyClient.R().SetResult(&tfir).
		SetBody(data).
		Post(fmt.Sprintf("https://miniprogram-kyc.tencentcloudapi.com/api/server/getfaceid?orderNo=%s", orderNo))
	if err != nil {
		slog.Error("获取人脸认证token失败", "err", err, "requestId", ctx.RequestId)
		return nil, err
	}
	//slog.Info("", "result", resp.Result().(*TencentFaceIdResponse))
	if tfir.Code != "0" {
		slog.Error("访问腾讯云faceId接口失败", "err", err, "result", tfir, "requestId", ctx.RequestId)
		return nil, errors.New(tfir.String())
	}

	faceId := tfir.Result.FaceID

	return TencentFaceIdResult{
		FaceID:        faceId,
		AgreementNo:   orderNo,
		OpenAPINonce:  nonce,
		OpenAPIUserID: userId,
		OpenAPISign:   sign,
	}, nil
}

type TencentFaceResultResponse struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	BizSeqNo string `json:"bizSeqNo"`
	Result   struct {
		OrderNo      string `json:"orderNo"`
		LiveRate     string `json:"liveRate"`
		Similarity   string `json:"similarity"`
		OccurredTime string `json:"occurredTime"`
		AppID        string `json:"appId"`
		Photo        string `json:"photo"`
		Video        string `json:"video"`
		BizSeqNo     string `json:"bizSeqNo"`
		SdkVersion   string `json:"sdkVersion"`
		TrtcFlag     string `json:"trtcFlag"`
	} `json:"result"`
	TransactionTime string `json:"transactionTime"`
}

func (f TencentFaceResultResponse) String() string {
	res, _ := sonic.ConfigFastest.MarshalToString(f)
	return res
}

func (tf TencentFace) GetTencentFaceResult(ctx utils.Context, userId string, orderNo string) (similarity TencentFaceResultResponse, err error) {
	// 人脸核身结果查询
	nonce := utils.GenterWithoutRepetitionStr(32)
	sign, err := tf.cookTencentSign(ctx, "", orderNo, nonce)
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

	similarity = TencentFaceResultResponse{}
	_, err = tf.RestyClient.R().SetResult(&similarity).SetBody(data).Post(url)
	if err != nil {
		return
	}
	if similarity.Code != "0" {
		slog.Error("访问腾讯云人脸核身结果查询接口失败", "err", err, "result", similarity, "requestId", ctx.RequestId)
		return
	}

	return similarity, nil
}
