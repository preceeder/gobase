package tencentyun

import (
	"encoding/json"
	"github.com/preceeder/gobase/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	v20180301 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
	"log/slog"
	"os"
)

var FaceIdClient *v20180301.Client

func InitTencentFaceId(appid string, serverSecret string) {
	credential := common.NewCredential(
		os.Getenv(appid),
		os.Getenv(serverSecret))
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	cpf.HttpProfile.ReqMethod = "POST"

	FaceIdClient, _ = v20180301.NewClient(credential, regions.Shanghai, cpf)
}

func Identity(ctx utils.Context, name string, idCard string) string {

	//创建common client
	m := map[string]string{"IdCard": idCard, "name": name}
	mStr, _ := json.Marshal(m)
	request := v20180301.IdCardVerificationRequest{}
	_ = request.FromJsonString(string(mStr))
	cardVerification, err := FaceIdClient.IdCardVerification(&request)
	if err != nil {
		slog.Error("访问身份二要素检查接口失败", "errors", err.Error(), "requestId", ctx.RequestId)
		return ""
	}
	return *cardVerification.Response.Result
}
