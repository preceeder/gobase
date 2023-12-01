package tencentyun

import (
	"github.com/preceeder/gobase/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	v20180301 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
	"log/slog"
)

var FaceIdClient *v20180301.Client

func InitTencentFaceId(appid string, serverSecret string) {
	credential := common.NewCredential(
		appid,
		serverSecret)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "faceid.tencentcloudapi.com"
	cpf.HttpProfile.ReqMethod = "POST"

	FaceIdClient, _ = v20180301.NewClient(credential, regions.Shanghai, cpf)
}

func Identity(ctx utils.Context, name string, idCard string) (string, string) {

	//创建common client
	request := v20180301.NewIdCardVerificationRequest()
	request.IdCard = &idCard
	request.Name = &name
	cardVerification, err := FaceIdClient.IdCardVerification(request)
	if err != nil {
		slog.Error("访问身份二要素检查接口失败", "errors", err.Error(), "requestId", ctx.RequestId)
		return "", "服务异常"
	}
	return *cardVerification.Response.Result, *cardVerification.Response.Description
}
