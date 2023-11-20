/*
File Name:  mysqlType.go
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

type DateTime carbon.DateTime

//	func (d Date) Value() (driver.Value, error) {
//		tempTime, _ := time.Parse("%Y-%m-%d", string(d))
//		return driver.Value(string(tempTime.UnixMilli())), nil
//	}
func (g *DateTime) Scan(src interface{}) error {
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
	v := carbon.ParseByFormat(string(source), "2006-01-02 15:04:05")
	*g = DateTime(carbon.DateTime{Carbon: v})
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

type NullString string

func (i *NullString) Scan(src any) error {
	if src == nil {
		*i = ""
		return nil
	}
	err := convertAssign(i, src)
	return err
}

type NullInt int

func (i *NullInt) Scan(src any) error {
	if src == nil {
		*i = 0
		return nil
	}
	err := convertAssign(i, src)
	return err
}

type NullInt8 int8

func (i *NullInt8) Scan(src any) error {
	if src == nil {
		*i = 0
		return nil
	}
	err := convertAssign(i, src)
	return err
}

type NullInt16 int16

func (i *NullInt16) Scan(src any) error {
	if src == nil {
		*i = 0
		return nil
	}
	err := convertAssign(i, src)
	return err
}

type NullInt32 int32

func (i *NullInt32) Scan(src any) error {
	if src == nil {
		*i = 0
		return nil
	}
	err := convertAssign(i, src)

	return err
}

type NullInt64 int64

func (i *NullInt64) Scan(src any) error {
	if src == nil {
		*i = 0
		return nil
	}
	err := convertAssign(i, src)

	return err
}

type NullBool bool

func (i *NullBool) Scan(src any) error {
	if src == nil {
		*i = false
		return nil
	}
	err := convertAssign(i, src)
	return err
}

type NullByte byte

func (i *NullByte) Scan(src any) error {
	if src == nil {
		*i = 0
		return nil
	}
	err := convertAssign(i, src)
	return err
}

type NullFloat64 float64

func (i *NullFloat64) Scan(src any) error {
	if src == nil {
		*i = 0
		return nil
	}
	err := convertAssign(i, src)

	return err
}
