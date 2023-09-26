/*
File Name:  ios.py
Description:
Author:      Chenghu
Date:       2023/8/23 22:28
Change Activity:
*/
package pay

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"time"
)

const (
	IosSandboxUrl        = "https://sandbox.itunes.apple.com/verifyReceipt" // (沙盒)链接
	IosProductionUrl     = "https://buy.itunes.apple.com/verifyReceipt"     //  (正式)链接
	IosRequestUrlTimeout = 15 * time.Second                                 // 请求接口超时时间
)

var IosResponseStatusMsg = map[int64]string{
	21000: "App Store 的请求不是使用 HTTP POST 请求方法发出的。", //·21000 对 App Store 的请求不是使用 HTTP POST 请求方法发出的。
	21001: "App Store 不再发送此状态代码。",                 //·21001 App Store 不再发送此状态代码。
	21002: "属性中的数据receipt-data格式不正确或服务遇到临时问题。",    //·21002 属性中的数据receipt-data格式不正确或服务遇到临时问题。再试一次。
	21003: "无法验证收据。",                              //·21003 无法验证收据。
	21004: "您提供的共享密钥与您帐户的文件共享密钥不匹配。",              //·21004 您提供的共享密钥与您帐户的文件共享密钥不匹配。
	21005: "收据服务器暂时无法提供收据。",                       //·21005 收据服务器暂时无法提供收据。再试一次
	21006: "此收据有效，但订阅已过期。",                        //·21006 此收据有效，但订阅已过期。当此状态代码返回到您的服务器时，收据数据也会被解码并作为响应的一部分返回。仅针对自动续订订阅的 iOS 6 样式交易收据返回。
	21007: "这条回执是来自测试环境，但它是发送到生产环境进行验证的。",         //·21007 这条回执是来自测试环境，但它是发送到生产环境进行验证的。
	21008: "这条回执来自生产环境，但它被发送到测试环境进行验证。",           //·21008 这条回执来自生产环境，但它被发送到测试环境进行验证。
	21009: "内部数据访问错误。",                            //·21009 内部数据访问错误。稍后再试。
	21010: "用户帐户找不到或已被删除。",                        //·21010 用户帐户找不到或已被删除。
}

// 官方文档 https://developer.apple.com/documentation/appstorereceipts/verifyreceipt
type IosPay struct {
	Params   IosParamsPay   `json:"params"`   //请求参数
	Response IosResponsePay `json:"response"` // 返回参数
	Error    error          `json:"error"`    //错误信息
}

func NewIosPay() *IosPay {
	return &IosPay{
		Params:   IosParamsPay{},
		Response: IosResponsePay{},
		Error:    nil,
	}
}

// 请求参数
type IosParamsPay struct {
	AppleClientId string `json:"appleClientId"` // 收据所属应用的捆绑包标识符(app应用的bundle_id)
	Password      string `json:"password"`      // 您的应用程序的共享密钥，它是一个十六进制字符串   有的话就填上
	Url           string `json:"url"`           // 请求链接    内部自动使用默认的
	//TransactionId      string `json:"transactionId"`      // 交易的唯一标识符
	TransactionReceipt string `json:"transactionReceipt"` //Base64 编码的收据数据
}
type InApp struct {
	CancellationDate      string `json:"cancellation_date"`       // App Store 退还交易或从家庭共享中撤销交易的时间，采用类似于 ISO 8601 的日期时间格式。此字段仅适用于已退款或撤销的交易。
	CancellationDateMs    string `json:"cancellation_date_ms"`    // App Store 退还交易或从家庭共享中撤销交易的时间，采用 UNIX 纪元时间格式，以毫秒为单位。此字段仅适用于退款或撤销的交易。使用此时间格式处理
	CancellationDatePst   string `json:"cancellation_date_pst"`   // App Store 退还交易或从家庭共享中撤销交易的时间，在太平洋时区。此字段仅适用于退款或撤销的交易。
	CancellationReason    string `json:"cancellation_reason"`     // 退款或撤销交易的原因。值“1”表示客户由于您的应用程序中的实际或感知问题而取消了他们的交易。值“0”表示交易因其他原因被取消；例如，如果客户意外购买。  可能的值：1, 0
	ExpiresDate           string `json:"expires_date"`            // 订阅到期时间或续订时间，采用类似于 ISO 8601 的日期时间格式。
	ExpiresDateMs         string `json:"expires_date_ms"`         // 订阅到期或续订的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式处理日期
	ExpiresDatePst        string `json:"expires_date_pst"`        // 太平洋时区的订阅到期时间或续订时间。
	OriginalTransactionId string `json:"original_transaction_id"` // 原始购买的交易标识符。
	ProductId             string `json:"product_id"`              // 购买的产品的唯一标识符。您在 App Store Connect 中创建产品时提供此值，它对应于存储在交易的支付属性中的对象的属性
	PromotionalOfferId    string `json:"promotional_offer_id"`    // 用户兑换的订阅优惠的标识符
	Quantity              string `json:"quantity"`                // 购买的消耗品数量。此值对应于SKPayment存储在交易的支付属性中的对象的数量属性。“1”除非使用可变付款进行修改，否则该值通常是不变的。最大值为 10。
	TransactionId         string `json:"transaction_id"`          // 交易的唯一标识符，例如购买、恢复或续订
	WebOrderLineItemId    string `json:"web_order_line_item_id"`  // 跨设备购买事件的唯一标识符，包括订阅续订事件。该值是识别订阅购买的主键。
}

// 内购支付返回的结构体 https://developer.apple.com/documentation/appstorereceipts/responsebody
type IosResponsePay struct {
	Status      *int64 `json:"status"`      //具体注释看 config.ResponseStatusMsg 0:正确,其他失败
	Environment string `json:"environment"` //（生成收据的环境） 值：Sandbox(沙盒), Production(正式环境)
	Receipt     struct {
		ReceiptType                string  `json:"receipt_type"`                // 生成的收据类型。该值对应于进行应用程序或 VPP 购买的环境。可能的值：Production, ProductionVPP, ProductionSandbox, ProductionVPPSandbox
		AdamId                     int64   `json:"adam_id"`                     //见。app_item_id
		AppItemId                  int64   `json:"app_item_id"`                 // 由 App Store Connect 生成并由 App Store 用于唯一标识所购买的应用程序。仅在生产中为应用分配此标识符。将此值视为 64 位长整数。
		ApplicationVersion         string  `json:"application_version"`         // 应用程序的版本号。应用程序的版本号对应于. 在生产中，此值是设备上基于. 在沙盒中，该值始终为。CFBundleVersionCFBundleShortVersionStringInfo.plistreceipt_creation_date_ms"1.0"
		BundleId                   string  `json:"bundle_id"`                   //收据所属应用的捆绑包标识符
		DownloadId                 int64   `json:"download_id"`                 //应用下载交易的唯一标识符。
		VersionExternalIdentifier  int64   `json:"version_external_identifier"` // 一个任意数字，用于标识您的应用程序的修订版。在沙箱中，此键的值为“0”。
		ReceiptCreationDate        string  `json:"receipt_creation_date"`       // App Store 生成收据的时间，采用类似于 ISO 8601 的日期时间格式。
		ReceiptCreationDateMs      string  `json:"receipt_creation_date_ms"`    // App Store 生成收据的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式处理日期。这个值不会改变。
		ReceiptCreationDatePst     string  `json:"receipt_creation_date_pst"`   // App Store 生成收据的时间，在太平洋时区。
		RequestDate                string  `json:"request_date"`                // 处理对端点的请求并生成响应的时间，采用类似于 ISO 8601 的日期时间格式。verifyReceipt
		RequestDateMs              string  `json:"request_date_ms"`             // 处理对端点的请求并生成响应的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式处理日期。verifyReceipt
		RequestDatePst             string  `json:"request_date_pst"`            // 在太平洋时区处理对端点的请求并生成响应的时间。verifyReceipt
		OriginalPurchaseDate       string  `json:"original_purchase_date"`      // 原始应用购买的时间，采用类似于 ISO 8601 的日期时间格式。
		OriginalPurchaseDateMs     string  `json:"original_purchase_date_ms"`   // 原始应用购买的时间，采用 UNIX 纪元时间格式，以毫秒为单位。使用此时间格式处理日期。
		OriginalPurchaseDatePst    string  `json:"original_purchase_date_pst"`  // 原始应用购买时间，太平洋时区。
		OriginalApplicationVersion string  `json:"original_application_version"`
		PreorderDate               string  `json:"preorder_date"`       // 用户订购可预购应用的时间，采用类似于 ISO 8601 的日期时间格式。
		PreorderDateMs             string  `json:"preorder_date_ms"`    // 用户订购可用于预订的应用程序的时间，采用 UNIX 纪元时间格式，以毫秒为单位。此字段仅在用户预购应用程序时出现。使用此时间格式处理日期。
		PreorderDatePst            string  `json:"preorder_date_pst"`   // 用户在太平洋时区订购可供预订的应用程序的时间。
		ExpirationDate             string  `json:"expiration_date"`     // 通过批量购买计划购买的应用程序ray1234据到期时间，采用类似于 ISO 8601 的日期时间格式。
		ExpirationDateMs           string  `json:"expiration_date_ms"`  // 通过批量购买计划购买的应用程序的收据到期时间，采用 UNIX 纪元时间格式，以毫秒为单位。如果通过批量购买计划购买的应用程序没有此密钥，则收据不会过期。使用此时间格式处理日期。
		ExpirationDatePst          string  `json:"expiration_date_pst"` // 通过批量购买计划购买的应用程序在太平洋时区的收据到期时间。
		InApp                      []InApp `json:"in_app"`              //包含所有应用内购买交易的应用内购买收据字段的数组
	} `json:"receipt"` //发送以供验证的收据的 JSON 表示形式。
}

func (i *IosPay) Send(url string, params map[string]string, resData any) error {
	//重试条件
	f := func(res *resty.Response, err error) bool {
		//重试次数设置为 2 测这个方法会执行 3次
		if err != nil {
			return true
		}
		return false
	}
	headers := map[string]string{"Content-Type": "application/json", "Accept": "application/json"}
	_, err := resty.New().AddRetryCondition(f).SetRetryCount(2).
		SetTimeout(IosRequestUrlTimeout).R().
		SetHeaders(headers).
		ForceContentType("application/json").
		SetBody(params).
		EnableTrace().
		SetResult(resData).
		Post(url)

	return err

}

// 请求苹果内购校验收据接口
func (this *IosPay) IosRequestReceiptUrl() *IosPay {
	if this.Error != nil {
		return this
	}
	if this.Params.TransactionReceipt == "" {
		this.Error = errors.New("缺少凭证数据")
		return this
	}
	// 请求数据
	params := map[string]string{
		"receipt-data": this.Params.TransactionReceipt,
		//ExcludeOldTransactions bool   `json:"exclude-old-transactions"` // 将此值设置true为以使响应仅包含任何订阅的最新续订交易。此字段仅用于包含自动续订订阅的应用收据。
	}
	if this.Params.Password != "" {
		params["password"] = this.Params.Password
	}
	resData := &IosResponsePay{}
	err := this.Send(this.Params.Url, params, resData)
	this.Response = *resData
	if err != nil {
		this.Error = errors.New(fmt.Sprintf("校验超时, 请稍后再试:%v", err))
		return this
	}

	return this
}

// 正式环境-参数
func (this *IosPay) AppleProductionParams(params IosParamsPay) *IosPay {
	this.Params = params
	this.Params.Url = IosProductionUrl
	return this
}

// 沙盒环境-参数
func (this *IosPay) AppleSandboxParams(params IosParamsPay) *IosPay {
	this.Params = params
	this.Params.Url = IosSandboxUrl
	return this
}

// 自定义校验
func (this *IosPay) AppleVerify() *IosPay {
	if this.Error != nil {
		return this
	}
	result := this.Response
	if result.Status == nil {
		this.Error = errors.New("url request fail")
		return this
	}
	//slog.Info("apple info", "data", result)

	if *result.Status != 0 {
		if statusMsg, ok := IosResponseStatusMsg[*result.Status]; ok {
			this.Error = errors.New(statusMsg)
			return this
		}
		this.Error = errors.New(fmt.Sprintf("返回的状态异常:%d", result.Status))
		return this
	}
	// result.Receipt 发送以供验证的收据的 JSON 表示形式。
	if len(result.Receipt.InApp) <= 0 {
		this.Error = errors.New("空凭证")
		return this
	}
	if result.Receipt.BundleId != this.Params.AppleClientId {
		this.Error = errors.New("bundleID 无效")
		return this
	}
	// --- > 业务处理
	//var (
	//	rawTransactionId = ""
	//	rawProductId     = ""
	//)
	// 取出对应凭证
	//for i := 0; i < len(result.Receipt.InApp); i++ {
	//	val := result.Receipt.InApp[i]
	//	if this.Params.TransactionId == val.TransactionId {
	//		rawProductId = val.ProductId
	//		rawTransactionId = val.TransactionId
	//		break
	//	}
	//}
	//
	//if len(rawTransactionId) <= 0 {
	//	this.Error = errors.New("raw内无对应凭证")
	//	return this
	//}
	//if len(rawProductId) <= 0 {
	//	this.Error = errors.New("raw内无对应产品")
	//	return this
	//}

	return this
}

// 公共校验(优先校验正式环境,再校验沙盒环境)
func AppleCommonVerify(params IosParamsPay) *IosPay {
	response := NewIosPay().AppleProductionParams(params).IosRequestReceiptUrl().AppleVerify()
	if response.Error != nil {
		slog.Error("ios product request error", "error", response.Error.Error())
		return NewIosPay().AppleSandboxParams(params).IosRequestReceiptUrl().AppleVerify()
	}
	return response
}
