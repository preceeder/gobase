/*
File Name:  SliceHandler.go
Description:
Author:      Chenghu
Date:       2023/8/21 10:07
Change Activity:
*/
package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

// 简单数组 转化为字符串
func SliceToString[T int | string | float64 | float32 | byte](arr []T) (result string) {
	var tempString []string
	for _, i := range arr { //遍历数组中所有元素追加成string
		tempString = append(tempString, fmt.Sprint(i))
	}
	result = strings.Join(tempString, ",")
	return
}

// SliceToMap
// two slice  to map
func SliceToMap[K ~string | ~int, V any | ~string | ~int | ~uint8](k []K, v []V) (map[K]V, error) {
	if len(k) > len(v) {
		return nil, errors.New("k len lt v len")
	}
	var res = make(map[K]V)
	for i, k := range k {
		res[k] = v[i]
	}
	return res, nil
}

type CanHashType interface {
	~string | ~int | ~uint8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// 对象数组 转化为 对象 map
func StructSliceToStructMap[T ~[]E, E any, K CanHashType](datas T, f func(data E) (K, E)) map[K]E {
	var temp = map[K]E{}
	for _, v := range datas {
		k, value := f(v)
		temp[k] = value
	}
	return temp
}

// map 数组 转化为 strcut 数组  []map[string]any   to []struct{}
func LiceMapToStruct(rows []map[string]any, dest interface{}) error {
	var vp reflect.Value

	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	direct := reflect.Indirect(value)

	//目标必须是 对象数组
	slice, err := IsTargetType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}
	direct.SetLen(0)

	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := ReflectBaseType(slice.Elem())

	for _, mp := range rows {
		vp = reflect.New(base)
		vpl := reflect.Indirect(vp)
		for key, v := range mp {
			desTp := vpl.FieldByName(key).Type()
			vpl.FieldByName(key).Set(reflect.ValueOf(v).Convert(desTp))
		}
		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}
	}
	fmt.Printf("%#v\n", dest)
	return nil
}

// any类型的切片 转化为 string 类型的切片
func SliceAnyToSliceString(data []any) []string {
	ss := make([]string, len(data))
	for i, iface := range data {
		val := toString(iface)
		ss[i] = val
	}
	return ss
}

func toString(val interface{}) string {
	switch val := val.(type) {
	case string:
		return val
	default:
		return ""
	}
}

// 特定类型的slice  转化为 any类型，   any 类型转化为 特定类型使用  cast.ToStringSlice()  就可以
func SliceConvertToAny[P any | string | int](src []P) []any {
	var dest = make([]any, len(src))
	for i, v := range src {
		dest[i] = v
	}
	return dest
}
