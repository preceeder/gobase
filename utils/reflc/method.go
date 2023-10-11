/*
File Name:  method.py
Description:
Author:      Chenghu
Date:       2023/10/11 16:39
Change Activity:
*/
package reflc

import (
	"reflect"
)

type Method struct {
	Servers map[string]reflect.Method
	Rcvr    reflect.Value
	Typ     reflect.Type
}

// rep 的 传入 指针 非指针都可以
func MakeService(rep interface{}) *Method {
	ser := Method{}
	ser.Typ = reflect.TypeOf(rep)
	ser.Rcvr = reflect.ValueOf(rep)
	//name := reflect.Indirect(ser.Rcvr).Type().Name()
	ser.Servers = map[string]reflect.Method{}
	for i := 0; i < ser.Typ.NumMethod(); i++ {
		method := ser.Typ.Method(i)
		mname := method.Name // string
		ser.Servers[mname] = method
	}

	return &ser
}
