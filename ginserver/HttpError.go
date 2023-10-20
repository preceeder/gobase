//   File Name:  HttpError.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/20 09:32
//    Change Activity:

package ginserver

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
