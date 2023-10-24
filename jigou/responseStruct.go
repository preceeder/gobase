/*
File Name:  responseStruct.go
Description:
Author:      Chenghu
Date:       2023/8/22 17:14
Change Activity:
*/
package jigou

type UserCountList struct {
	RoomID         string `json:"RoomId"`
	UserCount      int    `json:"UserCount"`
	AdminUserCount int    `json:"AdminUserCount"`
}
type DataUserCountList struct {
	UserCountList []UserCountList `json:"UserCountList"`
}

type RoomNumbers struct {
	Code      int               `json:"Code"`
	Message   string            `json:"Message"`
	RequestID string            `json:"RequestId"`
	Data      DataUserCountList `json:"Data"`
}

type PublicResponse struct {
	Code      int    `json:"Code"`
	Message   string `json:"Message"`
	RequestID string `json:"RequestId"`
}

type SendCustomCommand struct {
	Code      int           `json:"Code"`
	Message   string        `json:"Message"`
	RequestID string        `json:"RequestId"`
	Data      DataFailUsers `json:"Data"`
}
type FailUsers struct {
	UID  string `json:"Uid"`
	Code int    `json:"Code"`
}
type DataFailUsers struct {
	FailUsers []FailUsers `json:"FailUsers"`
}

type GenerateIdentifyToken struct {
	Code      int         `json:"Code"`
	Data      IdentiToken `json:"Data"`
	Message   string      `json:"Message"`
	RequestID string      `json:"RequestId"`
}
type IdentiToken struct {
	IdentifyToken string `json:"IdentifyToken"`
	RemainTime    int    `json:"RemainTime"`
}

// token业务扩展：权限认证属性
type RtcRoomPayLoad struct {
	RoomId       string      `json:"room_id"`        //房间 id（必填）；用于对接口的房间 id 进行强验证
	Privilege    map[int]int `json:"privilege"`      //权限位开关列表；用于对接口的操作权限进行强验证
	StreamIdList []string    `json:"stream_id_list"` //流列表；用于对接口的流 id 进行强验证；允许为空，如果为空，则不对流 id 验证
}
