//   File Name:  common.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/6 10:42
//    Change Activity:

package tencentIm

import (
	"github.com/mitchellh/mapstructure"
	"log/slog"
)

// im 域名
var TencentImHost = "https://console.tim.qq.com"

// im 部分常用接口
var ApiMap = map[string]string{
	"AccountImport":      "/v4/im_open_login_svc/account_import", // 导入单个账号
	"MultiaccountImport": "/v4/im_open_login_svc/multiaccount_import",
	"AccountDelete":      "/v4/im_open_login_svc/account_delete", // 删除账号
	"AccountCheck":       "/v4/im_open_login_svc/account_check",  // 查询账号

	"SendMsg":      "/v4/openim/sendmsg",           // 发送单聊消息
	"BatchSendMsg": "/v4/openim/batchsendmsg",      // 批量发送单聊消息
	"QueryMsg":     "/v4/openim/admin_getroammsg",  // 查询单聊消息
	"MsgWithdraw":  "/v4/openim/admin_msgwithdraw", // 撤回单聊消息

	"ImPush":  "/v4/all_member_push/im_push",     //全员推送
	"GetAttr": "/v4/all_member_push/im_get_attr", // 获取用户属性
	"SetAttr": "/v4/all_member_push/im_set_attr", // 设置用户属性

	"PortraitGet": "/v4/profile/portrait_get", // 资料获取
	"PortraitSet": "/v4/profile/portrait_set", // 资料设置

	"GetRecentContact":    "/v4/recentcontact/get_list", // 拉取回话列表
	"DeleteRecentContact": "/v4/recentcontact/delete",   // 删除单个会话
}

type BaseResponse interface {
	GetResponse() map[string]any
}

type CommonResponse struct {
	ActionStatus string `json:"ActionStatus" mapstructure:"ActionStatus"`
	ErrorInfo    string `json:"ErrorInfo" mapstructure:"ErrorInfo"`
	ErrorCode    int    `json:"ErrorCode" mapstructure:"ErrorCode"`
}

func (b CommonResponse) GetResponse() (data map[string]any) {
	err := mapstructure.Decode(b, &data)
	if err != nil {
		slog.Error("获取tencent im 返回数据错误", "error", err.Error())
	}
	return
}
