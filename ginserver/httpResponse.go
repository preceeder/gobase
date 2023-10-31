//   File Name:  httpResponse.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/31 11:10
//    Change Activity:

package ginserver

import "reflect"

type HttpResponse interface {
	GetResponse() any
}

var HttpResponseType = reflect.TypeOf((*HttpResponse)(nil)).Elem()

type BaseHttpResponse struct {
	Success bool `json:"success"`
	Code    bool `json:"code"`
	Data    any  `json:"data"`
}

func (h BaseHttpResponse) GetResponse() any {
	return h
}
