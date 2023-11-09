//   File Name:  im.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/6 10:42
//    Change Activity:

package tencentIm

import (
	"github.com/go-resty/resty/v2"
	"github.com/preceeder/gobase/tencentIm/ECDSASHA256"
	"github.com/preceeder/gobase/tencentIm/HMACSHA256"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var ImClient *resty.Client
var TencentImConfig Config
var UserSign string // 生成的有效 usersign

type Config struct {
	AppId      int    `json:"appId"`      // appid
	Identifier string `json:"identifier"` // 管理员账户
	Key        string `json:"Key"`        // 密钥   HMAC-SHA256 算法 使用
	PrivateKey string `json:"privateKey"` // 私钥   ECDSA-SHA256 算法 加密 使用
	PublicKey  string `json:"publicKey"`  // 公钥    ECDSA-SHA256 算法 验证 使用
	UseSha     string `json:"useSha"`     // 使用那种算法  HMAC-SHA256｜ ECDSA-SHA256
	ImHost     string `json:"imHost"`     // 域名  最后不要 /
	Expire     int    `json:"expire"`     // token 过期时间 s
}

func InitWithViper(config viper.Viper) {
	utils.ReadViperConfig(config, "tencentIm", &TencentImConfig)
	InitIm()
}

func InitWithStruct() {
	InitIm()
}

func InitIm() {
	ImClient = resty.New()
	ImClient.SetTimeout(3 * time.Second)
	ImClient.SetTransport(&http.Transport{
		MaxIdleConnsPerHost:   50,               // 对于每个主机，保持最大空闲连接数为 10
		IdleConnTimeout:       30 * time.Second, // 空闲连接超时时间为 30 秒
		TLSHandshakeTimeout:   3 * time.Second,  // TLS 握手超时时间为 10 秒
		ResponseHeaderTimeout: 3 * time.Second,  // 等待响应头的超时时间为 20 秒
	})
}

// respBody 必须是指针

func SendImRequest(ctx utils.Context, serverName string, requestData any, respBody any) error {
	req := ImClient.R().SetBody(requestData)
	if respBody != nil {
		req.SetResult(respBody)
	}
	durl := setUrl(serverName)
	_, err := req.Post(durl)
	if err != nil {
		slog.Error("SendImRequest error", "error", err.Error(), "serverName", serverName, "data", requestData, "requestId", ctx.RequestId)
		return err
	}
	return nil
}

func setUrl(serverName string) string {
	if uri, ok := ApiMap[serverName]; ok {
		query := url.Values{}
		query.Set("contenttype", "json")
		query.Set("sdkappid", strconv.Itoa(TencentImConfig.AppId))
		query.Set("identifier", TencentImConfig.Identifier)
		sign := ""
		if TencentImConfig.UseSha == "HMAC-SHA256" {
			sign = getHmacSign()
		} else {
			sign = getEcdsaSign()
		}
		query.Set("usersig", sign)
		query.Set("random", utils.RandStrInt(5))
		bp, _ := url.JoinPath(TencentImConfig.ImHost, uri)
		return bp + "?" + query.Encode()
	} else {
		slog.Error("im server not find in ApiMap", "serverName", serverName)
	}
	return ""
}

func getHmacSign() string {
	var err error
	var userSignValid = false
	if len(UserSign) > 10 {
		err = HMACSHA256.VerifyUserSig(uint64(TencentImConfig.AppId), TencentImConfig.Key, TencentImConfig.Identifier, UserSign, time.Now())
		if err != nil {
			slog.Error("im usersign error", "error", err.Error())
		}
		userSignValid = true
	}
	if !userSignValid {
		UserSign, err = HMACSHA256.GenUserSig(TencentImConfig.AppId, TencentImConfig.Key, TencentImConfig.Identifier, TencentImConfig.Expire)
		if err != nil {
			slog.Error("生成im usersign error", "error", err.Error())
		}
	}
	return UserSign
}

func getEcdsaSign() string {
	var err error
	var userSignValid = false
	if len(UserSign) > 10 {
		err = ECDSASHA256.VerifyUsersig(TencentImConfig.PublicKey, UserSign, TencentImConfig.AppId, TencentImConfig.Identifier)
		if err != nil {
			slog.Error("im usersign error", "error", err.Error())
		}
		userSignValid = true
	}
	if !userSignValid {
		UserSign, err = ECDSASHA256.GenerateUsersigWithExpire(TencentImConfig.PrivateKey, TencentImConfig.AppId, TencentImConfig.Identifier, int64(TencentImConfig.Expire))
		if err != nil {
			slog.Error("生成im usersign error", "error", err.Error())
		}
	}
	return UserSign
}

func GetUserSign(userId string) (userSign string, err error) {
	if TencentImConfig.UseSha == "ECDSA-SHA256" {
		userSign, err = ECDSASHA256.GenerateUsersigWithExpire(TencentImConfig.PrivateKey, TencentImConfig.AppId, userId, int64(TencentImConfig.Expire))
		if err != nil {
			slog.Error("生成im usersign error", "error", err.Error())
		}
	} else if TencentImConfig.UseSha == "HMAC-SHA256" {
		userSign, err = HMACSHA256.GenUserSig(TencentImConfig.AppId, TencentImConfig.Key, userId, TencentImConfig.Expire)
		if err != nil {
			slog.Error("生成im usersign error", "error", err.Error())
		}
	}
	return
}
