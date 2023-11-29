package face

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	facebody "github.com/alibabacloud-go/facebody-20191230/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
)

type ALFaceConfig struct {
	Name     string `json:"name"`
	KeyId    string `json:"keyId"`
	Secret   string `json:"secret"`
	EndPoint string `json:"endPoint"`
	RegionId string `json:"regionId"`
	AppKey   string `json:"appKey"`
	Env      string `json:"env"`
}

type ALFaceClientStruct struct {
	Client *facebody.Client
	Config ALFaceConfig
}

var AliFaceClient ALFaceClientStruct

func InitWithViper(config viper.Viper) {
	//aliConfig := readAliPushConfig(config)
	cnf := ALFaceConfig{}
	utils.ReadViperConfig(config, "ali_face", &cnf)
	client, err := CreateClient(cnf.KeyId, cnf.Secret, cnf.EndPoint, cnf.RegionId)
	if err != nil {
		return
	}
	if err != nil {
		slog.Error("阿里云人脸识别创建失败", "error", err.Error())
		panic("阿里云人脸识别创建失败：" + err.Error())
	}
	AliFaceClient = ALFaceClientStruct{Client: client, Config: cnf}
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId string, accessKeySecret string, endpoint string, regionId string) (_result *facebody.Client, _err error) {
	config := new(openapi.Config)

	// init config with ak
	config.SetAccessKeyId(accessKeyId).
		SetAccessKeySecret(accessKeySecret).
		SetRegionId(regionId).
		SetEndpoint(endpoint)

	// init config with credential
	credentialConfig := &credential.Config{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
		SecurityToken:   config.SecurityToken,
		Type:            tea.String("access_key"),
	}
	// If you have any questions, please refer to it https://github.com/aliyun/credentials-go/blob/master/README-CN.md
	cred, err := credential.NewCredential(credentialConfig)
	if err != nil {
		panic(err)
	}
	config.SetCredential(cred)

	// init client
	client, err := facebody.NewClient(config)
	if err != nil {
		panic(err)
	}
	return client, _err
}

func (alfc ALFaceClientStruct) CompareFace(ctx utils.Context, imageUrlA string, imageUrlB string) *facebody.CompareFaceResponse {
	httpClient := http.Client{}
	file1, _ := httpClient.Get(imageUrlA)
	file2, _ := httpClient.Get(imageUrlB)
	compareFaceRequest := &facebody.CompareFaceAdvanceRequest{
		ImageURLAObject: file1.Body,
		ImageURLBObject: file2.Body,
	}
	runtime := &util.RuntimeOptions{}
	compareFaceResponse, err := alfc.Client.CompareFaceAdvance(compareFaceRequest, runtime)
	if err != nil {
		// 获取整体报错信息
		slog.Error("人脸比较接口访问失败", "errors", err.Error(), "requestId", ctx.RequestId)
		return nil
	} else {
		// 获取整体结果
		return compareFaceResponse
	}
}

func (alfc ALFaceClientStruct) RecognizeFac(ctx utils.Context, imageUrl string) *facebody.RecognizeFaceResponse {
	httpClient := http.Client{}
	file, _ := httpClient.Get(imageUrl)
	recognizeFaceRequest := &facebody.RecognizeFaceAdvanceRequest{
		ImageURLObject: file.Body,
	}
	runtime := &util.RuntimeOptions{}
	recognizeFaceResponse, err := alfc.Client.RecognizeFaceAdvance(recognizeFaceRequest, runtime)
	if err != nil {
		// 获取整体报错信息
		slog.Error("人脸属性识别接口访问失败", "errors", err.Error(), "requestId", ctx.RequestId)
		return nil
	} else {
		// 获取整体结果
		return recognizeFaceResponse
	}
}
