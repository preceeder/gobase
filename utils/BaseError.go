/*
File Name:  HttpError.py
Description:
Author:      Chenghu
Date:       2023/8/29 14:14
Change Activity:
*/
package utils

type BaseHttpError struct {
	ErrorCode int
	Message   string
}

func (h BaseHttpError) GetMap() map[string]any {
	return map[string]any{"errorCode": h.ErrorCode, "message": h.Message}
}

func (h BaseHttpError) Error() string {
	return h.Message
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
