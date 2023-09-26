/*
File Name:  MapHandler.py
Description:
Author:      Chenghu
Date:       2023/8/21 10:48
Change Activity:
*/
package utils

import (
	"github.com/mitchellh/mapstructure"
	"log/slog"
)

// map 转为 strcut 同时还获取 meta数据  可以知道对应的 strcut 那些属性是没有的
func MapToStructWithMeta(input any, output any) (md mapstructure.Metadata) {
	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   output,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		slog.Error("mapToStruct mapstructure.NewDecoder(config)", "error", err.Error())
	}
	if err := decoder.Decode(input); err != nil {
		slog.Error("mapToStruct decoder.Decode(input)", "error", err.Error())

		panic(err)
	}
	return
}
