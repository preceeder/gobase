/*
File Name:  mysqlType.py
Description:
Author:      Chenghu
Date:       2023/8/16 11:16
Change Activity:
*/
package mysqlDb

import (
	"github.com/bytedance/sonic"
	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
)

// 将date  类型转化为 int 类型
type Date int

//	func (d Date) Value() (driver.Value, error) {
//		tempTime, _ := time.Parse("%Y-%m-%d", string(d))
//		return driver.Value(string(tempTime.UnixMilli())), nil
//	}
func (g *Date) Scan(src interface{}) error {
	var source []byte
	// let's support string and []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Incompatible type for GzippedText")
	}
	//v, _ := time.Parse("%Y-%m-%d", string(source))
	v := carbon.ParseByFormat(string(source), "Y-m-d")
	*g = Date(v.Timestamp())
	return nil
}

type Json map[string]any

func (j *Json) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Incompatible type for string")
	}
	err := sonic.Unmarshal(source, j)
	if err != nil {
		return errors.New("Incompatible type for string sonic.Unmarshal error: " + err.Error())
	}
	return nil
}
