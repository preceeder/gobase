/*
File Name:  HttpError.py
Description:
Author:      Chenghu
Date:       2023/8/29 14:14
Change Activity:
*/
package utils

import "github.com/bytedance/sonic"

type Error interface {
	GetData() string
	Error() string
}

type AllError struct {
	Message string
	Data    any
}

func (a AllError) GetData() string {
	data, _ := sonic.MarshalString(a.Data)
	return data
}

func (a AllError) Error() string {
	data, _ := sonic.MarshalString(map[string]string{
		"data":    a.GetData(),
		"message": a.Message,
	})
	return data
}

type BaseHttpError struct {
	Code      int
	ErrorCode int
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
	return h.Code
}

type BaseRunTimeError struct {
	ErrorCode int
	Message   string
}

func (h BaseRunTimeError) GetMap() map[string]any {
	return map[string]any{"errorCode": h.ErrorCode, "message": h.Message}
}

func (h BaseRunTimeError) Error() string {
	return h.Message
}
