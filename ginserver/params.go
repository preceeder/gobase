/*
File Name:  params.go
Description:
Author:      Chenghu
Date:       2023/10/12 06:26
Change Activity:
*/
package ginserver

import (
	"github.com/bytedance/sonic"
	"reflect"
)

type GinParam interface {
	GetType() string
	String() string
}

var GinParamType = reflect.TypeOf((*GinParam)(nil)).Elem()

// query
//type QueryString string
//
//func (q QueryString) GetType() string {
//	return "query"
//}
//
//type QueryInt int
//
//func (q QueryInt) GetType() string {
//	return "query"
//}
//
//type QueryFloat float64
//
//func (q QueryFloat) GetType() string {
//	return "query"
//}

// 结构体中需要 有tag -> `form:"data"`
type Query struct {
}

func (q Query) GetType() string {
	return "query"
}
func (receiver Query) String() string {
	data, _ := sonic.MarshalString(receiver)
	return data
}

type Header struct {
}

func (receiver Header) GetType() string {
	return "header"
}
func (receiver Header) String() string {
	data, _ := sonic.MarshalString(receiver)
	return data
}

// body json
type BodyJson struct {
}

func (receiver BodyJson) GetType() string {
	return "json"
}

func (receiver BodyJson) String() string {
	data, _ := sonic.MarshalString(receiver)
	return data
}

// form
type Form struct {
}

func (f Form) GetType() string {
	return "form"
}

func (receiver Form) String() string {
	data, _ := sonic.MarshalString(receiver)
	return data
}
