//   File Name:  code.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/19 10:15
//    Change Activity:

package ginserver

type ErrorCode int

const (
	CodeSystemError       ErrorCode = 10000 // 系统错误
	CodeParameterError    ErrorCode = 10001 // 参数错误
	CodeRequestFrequently ErrorCode = 10002 // 请求超限
	CodeLoginError        ErrorCode = 10003 // 登录失败
	CodeLoinForbid        ErrorCode = 10004 // 禁止登录
	CodeProhibit          ErrorCode = 10005 // 禁止操作 operation
)

var CodeMessage = map[ErrorCode]string{
	CodeSystemError:       "System error",
	CodeParameterError:    "Invalid parameter",
	CodeRequestFrequently: "Frequent requests",
	CodeLoginError:        "Login failed",
	CodeLoinForbid:        "Prohibit login",
	CodeProhibit:          "Prohibit operation",
}
