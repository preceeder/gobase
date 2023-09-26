/*
File Name:  rongyun.py
Description:
Author:      Chenghu
Date:       2023/8/23 21:44
Change Activity:
*/
package rongyun

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/rongcloud/server-sdk-go/v3/sdk"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var Rc *sdk.RongCloud
var rcConfig rongyunConfig

type rongyunConfig struct {
	AppKey           string `json:"appKey"`
	AppSecret        string `json:"appSecret"`
	WithRongCloudURI string `json:"withRongCloudURI"`
}

func InitWithViper(config viper.Viper) {
	ry := readRongYunConfig(config)
	rcConfig = ry
	InitRongYun(ry.AppKey, ry.AppSecret, ry.WithRongCloudURI)
}

func readRongYunConfig(v viper.Viper) (ry rongyunConfig) {
	rongyun := v.Sub("rongyun")
	if rongyun == nil {
		fmt.Printf("rongyun config is nil")
		os.Exit(1)
	}
	ry = rongyunConfig{}
	err := rongyun.Unmarshal(&ry)
	if err != nil {

		fmt.Printf("rongyun config read error: " + err.Error())
		os.Exit(1)
	}
	return
}

func InitRongYun(appKey string, appSecret string, withRongCloudURI string) {
	// 方法1： 创建对象时设置
	Rc = sdk.NewRongCloud(appKey,
		appSecret,
		// 每个域名最大活跃连接数
		sdk.WithMaxIdleConnsPerHost(100),
		sdk.WithTimeout(10),
		sdk.WithRongCloudURI(withRongCloudURI),
	)
}

var CustomMsgType = "CU:CUSTOM"

type CustomMsg struct {
	Class   string
	User    sdk.MsgUserInfo `json:"user"`
	Content string
}

func (c CustomMsg) ToString() (string, error) {
	msg, err := sonic.MarshalString(c)
	if err != nil {
		return "", err
	}
	return msg, nil
}

// 回调参数校验
func CallDataCheck(timestamp, nonce, signature string) bool {
	data := []string{rcConfig.AppSecret, nonce, timestamp}
	chd := cryptor.Sha1(strings.Join(data, ""))
	if chd == signature {
		return true
	}
	return false
}
