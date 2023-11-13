package face

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	facebody "github.com/alibabacloud-go/facebody-20191230/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
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

type ALFaceClient struct {
	Client *facebody.Client
	Config ALFaceConfig
}

var AliFaceClientList map[string]ALFaceClient = make(map[string]ALFaceClient)

func InitWithViper(config viper.Viper) {
	//aliConfig := readAliPushConfig(config)
	cnf := []ALFaceConfig{}
	utils.ReadViperConfig(config, "ali_face", &cnf)
	for _, cf := range cnf {
		client, err := CreateClient(&(cf.KeyId), &(cf.Secret), &(cf.EndPoint))
		if err != nil {
			return
		}
		if err != nil {
			slog.Error("阿里云人脸识别创建失败", "error", err.Error())
			panic("阿里云人脸识别创建失败：" + err.Error())
		}
		AliFaceClientList[cf.Name] = ALFaceClient{Client: client, Config: cf}
	}
}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string, endpoint *string) (_result *facebody.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = endpoint
	_result = &facebody.Client{}
	_result, _err = facebody.NewClient(config)
	return _result, _err
}

func (alfc ALFaceClient) CompareFace(ctx utils.Context, imageUrlA string, imageUrlB string) *facebody.CompareFaceResponse {
	compareFaceRequest := &facebody.CompareFaceRequest{
		ImageURLA: tea.String(imageUrlA),
		ImageURLB: tea.String(imageUrlB),
	}
	runtime := &util.RuntimeOptions{}
	compareFaceResponse, err := alfc.Client.CompareFaceWithOptions(compareFaceRequest, runtime)
	if err != nil {
		// 获取整体报错信息
		slog.Error("人脸比较接口访问失败", "errors", err.Error(), "requestId", ctx.RequestId)
		return nil
	} else {
		// 获取整体结果
		return compareFaceResponse
	}
}

func (alfc ALFaceClient) RecognizeFac(ctx utils.Context, imageUrl string) *facebody.RecognizeFaceResponse {
	recognizeFaceRequest := &facebody.RecognizeFaceRequest{
		ImageURL: tea.String(imageUrl),
	}
	runtime := &util.RuntimeOptions{}
	recognizeFaceResponse, err := alfc.Client.RecognizeFaceWithOptions(recognizeFaceRequest, runtime)
	if err != nil {
		// 获取整体报错信息
		slog.Error("人脸属性识别接口访问失败", "errors", err.Error(), "requestId", ctx.RequestId)
		return nil
	} else {
		// 获取整体结果
		return recognizeFaceResponse
	}
}
