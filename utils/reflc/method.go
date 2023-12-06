/*
File Name:  method.go
Description:
Author:      Chenghu
Date:       2023/10/11 16:39
Change Activity:
*/
package reflc

import (
	"fmt"
	"reflect"
	"strconv"
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

func Dushd(p reflect.Type, value string) reflect.Value {
	data := reflect.New(p)
	data.Set(reflect.ValueOf(value))
	return data
}

func Direct(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

// 传入指定的 类型， 我对应的字符串    使用场景 已知一个数据的 类型， 和一个已知的 字符串  且可以传化为那类型， 返回对应类型的reflect.value
func DUnmarshal(p reflect.Type, value string) (reflect.Value, error) {
	//去除前后的连续空格
	p = Direct(p)
	data := reflect.New(p)
	data = data.Elem() //解析指针

	switch p.Kind() {
	case reflect.String:
		data.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err == nil {
			data.SetBool(b)
		} else {
			return reflect.ValueOf(""), err
		}
	case reflect.Float32,
		reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.ValueOf(""), err
		} else {
			data.SetFloat(f) //通过reflect.Value修改原始数据的值
		}
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return reflect.ValueOf(""), err
		} else {
			data.SetInt(i) //有符号整型通过SetInt
		}
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.ValueOf(""), err
		} else {
			data.SetUint(i) //无符号整型需要通过SetUint
		}
	case reflect.Slice:
		//err := sonic.UnmarshalString(value, data.Interface())
		if value[0] == '[' && value[len(value)-1] == ']' {
			arr := SplitJson(value[1 : len(value)-1]) //去除前后的[]
			if len(arr) > 0 {
				data.Set(reflect.MakeSlice(p, len(arr), len(arr))) //通过反射创建slice
				for i := 0; i < len(arr); i++ {
					eleValue := data.Index(i)
					eleType := eleValue.Type()
					//if eleType.Kind() != reflect.Ptr {
					//	eleValue = eleValue.Addr()
					//}
					sliceD, err := DUnmarshal(eleType, arr[i])
					if err != nil {
						return reflect.ValueOf(""), err
					}
					eleValue.Set(sliceD)
				}
			}
		} else if value != "null" {
			return reflect.ValueOf(""), fmt.Errorf("invalid json part: %s", value)
		}
	default:
		fmt.Printf("暂不支持类型:%s\n", p.Kind().String())
	}
	return data, nil
}
