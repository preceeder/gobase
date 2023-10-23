//   File Name:  oss.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/18 14:58
//    Change Activity:

package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var AliOssConfig AliOss
var OssClient *oss.Client
var OssBucket *oss.Bucket // 使用的时候都是用这个

func InitAliOssWithViper(config viper.Viper) {
	//aliConfig := readAliOssConfig(config)
	utils.ReadViperConfig(config, "ali_oss", &AliOssConfig)
	CreateClient(AliOssConfig.EndPoint, AliOssConfig.AccessKey, AliOssConfig.AccessSecret)
}

type AliOss struct {
	AccessKey    string `json:"accessKey"`
	AccessSecret string `json:"accessSecret"`
	EndPoint     string `json:"endPoint"`
	Region       string `json:"region"`
	BucketName   string `json:"bucketName"`
	Cdn          string `json:"cdn"`
}

//func readAliOssConfig(v viper.Viper) (ali AliOss) {
//	alioss := v.Sub("ali_oss")
//	if alioss == nil {
//		fmt.Printf("ali oss config is nil")
//		os.Exit(1)
//	}
//	ali = AliOss{}
//	err := alioss.Unmarshal(&ali)
//	if err != nil {
//		fmt.Printf("ali oss config read error: " + err.Error())
//		os.Exit(1)
//	}
//	return
//}

func CreateClient(endpoint, accessKey, accessSecret string) {
	client, err := oss.New(endpoint, accessKey, accessSecret, oss.Region(AliOssConfig.Region))
	if err != nil {
		slog.Error("阿里云push创建失败", "error", err.Error())
		os.Exit(-1)
	}
	OssClient = client
	slog.Info("创建oss 客户端")
	bucket, err := OssClient.Bucket(AliOssConfig.BucketName)
	if err != nil {
		slog.Error("阿里云push创建失败", "error", err.Error())
		os.Exit(-1)
	}
	OssBucket = bucket
}
