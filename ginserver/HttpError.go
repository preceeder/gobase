//   File Name:  HttpError.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/20 09:32
//    Change Activity:

package ginserver

import "reflect"

type HttpError interface {
	GetCode() int // 正常情况都是 200， 错误情况一般是  403
	GetMap() map[string]any
	Error() string
}

var HttpErrorType = reflect.TypeOf((*HttpError)(nil)).Elem()

type BaseHttpError struct {
	Code      StatusCode
	ErrorCode ErrorCode
	Message   string
}

func (h BaseHttpError) GetMap() map[string]any {
	return map[string]any{"errorCode": h.ErrorCode, "message": h.Message}
}

func (h BaseHttpError) Error() string {
	return h.Message
}

func (h BaseHttpError) GetCode() int {
	if h.Code == 0 {
		return 403
	}
	return int(h.Code)
}
