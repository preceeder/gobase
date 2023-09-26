/*
File Name:  base.py
Description:
Author:      Chenghu
Date:       2023/8/21 10:08
Change Activity:
*/
package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

func ReflectBaseType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func IsTargetType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = ReflectBaseType(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

func RunFunc(object interface{}, methodName string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	//动态调用方法
	return reflect.ValueOf(object).MethodByName(methodName).Call(inputs)
}

func GetAttr(object interface{}, attr string, output interface{}) (err error) {
	//动态访问属性
	tv := reflect.TypeOf(object)
	var value interface{}
	if tv.Kind() == reflect.Ptr {
		value = reflect.ValueOf(object).Elem().FieldByName(attr).Interface()
	} else {
		value = reflect.ValueOf(object).FieldByName(attr).Interface()
	}
	err = mapstructure.Decode(value, &output)
	return
}

func AnyToString(agrs any, spacing []byte) (string, error) {
	//这里目前可以是 数组， int, string, float
	vV := reflect.ValueOf(agrs)
	vV = reflect.Indirect(vV)

	vT := reflect.TypeOf(agrs)
	if vT.Kind() == reflect.Ptr {
		vT = vT.Elem()
	}
	bf := bytes.Buffer{} //存放序列化结果

	switch vV.Kind() {
	case reflect.String:
		return fmt.Sprintf("%s", vV.String()), nil //取得reflect.Value对应的原始数据的值
	case reflect.Bool:
		return fmt.Sprintf("%t", vV.Bool()), nil
	case reflect.Float32,
		reflect.Float64:
		return fmt.Sprintf("%f", vV.Float()), nil
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return fmt.Sprintf("%v", vV.Interface()), nil
	case reflect.Slice:
		v := vT.Elem().Kind()
		fmt.Println(v)
		for i := 0; i < vV.Len(); i++ {
			if bs, err := AnyToString(vV.Index(i).Interface(), spacing); err != nil { //对slice的第i个元素进行序列化。递归
				return " ", err
			} else {
				bf.Write([]byte(bs))
				bf.Write(spacing)
			}
		}
	case reflect.Map:
		if vV.IsNil() {
			return " ", nil
		}
		if vV.Len() > 0 {
			for _, key := range vV.MapKeys() {
				if keyBs, err := AnyToString(key.Interface(), spacing); err != nil {
					return " ", err
				} else {
					bf.Write([]byte(keyBs))
					bf.WriteByte('=')
					v := vV.MapIndex(key)
					if vBs, err := AnyToString(v.Interface(), spacing); err != nil {
						return " ", err
					} else {
						bf.Write([]byte(vBs))
						bf.Write(spacing)
					}
				}
			}
			bf.Truncate(len(bf.Bytes()) - len(spacing)) //删除最后一个多余的spaceing
		}
	case reflect.Struct:
		//bf.WriteByte('{')
		if vV.NumField() > 0 {
			for i := 0; i < vV.NumField(); i++ {
				fieldValue := vV.Field(i)
				fieldType := vT.Field(i)
				name := fieldType.Name
				////如果没有json Tag，默认使用成员变量的名称
				//if len(fieldType.Tag.Get("json")) > 0 {
				//	name = fieldType.Tag.Get("json")
				//}
				bf.WriteString("\"")
				bf.WriteString(name)
				bf.WriteString("\"")
				bf.WriteString(":")
				if bs, err := AnyToString(fieldValue.Interface(), spacing); err != nil { //对value递归调用Marshal序列化
					return "", err
				} else {
					bf.Write([]byte(bs))
				}
				bf.Write(spacing)
			}
			bf.Truncate(len(bf.Bytes()) - len(spacing)) //删除最后一个逗号
		}
		//bf.WriteByte('}')
		//return bf.Bytes(), nil	default:
		//return "", errors.New("不支持的数据类型")
	}

	return bf.String(), nil
}

func AnyToSlice(agrs any) ([]any, error) {
	//这里目前可以是 数组， int, string, float
	if agrs == nil {
		panic("AnyToSlice can not temp")
	}
	vV := reflect.ValueOf(agrs)
	vV = reflect.Indirect(vV)

	vT := reflect.TypeOf(agrs)
	if vT.Kind() == reflect.Ptr {
		vT = vT.Elem()
	}
	//bf := bytes.Buffer{} //存放序列化结果
	bf := []any{}
	switch vV.Kind() {
	case reflect.String:
		bf = append(bf, vV.String())
		return bf, nil //取得reflect.Value对应的原始数据的值
	case reflect.Bool:
		bf = append(bf, vV.Bool())
		return bf, nil
	case reflect.Float32,
		reflect.Float64:
		bf = append(bf, vV.Float())
		return bf, nil
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		bf = append(bf, vV.Interface())
	case reflect.Slice:
		v := vT.Elem().Kind()
		fmt.Println(v)
		for i := 0; i < vV.Len(); i++ {
			if bs, err := AnyToSlice(vV.Index(i).Interface()); err != nil { //对slice的第i个元素进行序列化。递归
				return bf, err
			} else {
				bf = append(bf, bs...)
			}
		}
	default:
		return bf, errors.New("不支持的数据类型")
	}

	return bf, nil

}
