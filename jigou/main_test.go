/*
File Name:  main_test.go.py
Description:
Author:      Chenghu
Date:       2023/8/22 16:27
Change Activity:
*/
package jigou

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"net/url"
	"reflect"
	"testing"
)

func TestJiGou_GetRoomNumbers(t *testing.T) {
	type fields struct {
		AppId        string
		ServerSecret string
		RestyClient  *resty.Client
	}
	type args struct {
		roomId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{name: "jigou", fields: fields{
			AppId:        "",
			ServerSecret: "",
			RestyClient:  resty.New(),
		},
			args: args{
				roomId: "",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JiGou{
				AppId:        tt.fields.AppId,
				ServerSecret: tt.fields.ServerSecret,
				RestyClient:  tt.fields.RestyClient,
			}
			res, err := j.GetRoomNumbers(tt.args.roomId)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(res)
		})
	}
}

func TestCallDataCheck(t *testing.T) {
	type args struct {
		callbacksecret string
		timestamp      string
		nonce          string
		signature      string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "", args: args{signature: "",
			timestamp: "", nonce: "", callbacksecret: ""}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallDataCheck(tt.args.timestamp, tt.args.nonce, tt.args.signature); got != tt.want {
				t.Errorf("CallDataCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitJiGou(t *testing.T) {
	type args struct {
		appid        string
		serverSecret string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitJiGou(tt.args.appid, tt.args.serverSecret)
		})
	}
}

func TestInitWithViper(t *testing.T) {
	type args struct {
		config viper.Viper
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitWithViper(tt.args.config)
		})
	}
}

func TestJiGou_CloseRoom(t *testing.T) {
	type fields struct {
		AppId        string
		ServerSecret string
		RestyClient  *resty.Client
	}
	type args struct {
		RoomId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    PublicResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JiGou{
				AppId:        tt.fields.AppId,
				ServerSecret: tt.fields.ServerSecret,
				RestyClient:  tt.fields.RestyClient,
			}
			got, err := j.CloseRoom(tt.args.RoomId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloseRoom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloseRoom() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJiGou_GenerateIdentifyToken(t *testing.T) {
	type fields struct {
		AppId        string
		ServerSecret string
		RestyClient  *resty.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    GenerateIdentifyToken
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JiGou{
				AppId:        tt.fields.AppId,
				ServerSecret: tt.fields.ServerSecret,
				RestyClient:  tt.fields.RestyClient,
			}
			got, err := j.GenerateIdentifyToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateIdentifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateIdentifyToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJiGou_Get(t *testing.T) {
	type fields struct {
		AppId        string
		ServerSecret string
		RestyClient  *resty.Client
	}
	type args struct {
		url     string
		params  url.Values
		resBody any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JiGou{
				AppId:        tt.fields.AppId,
				ServerSecret: tt.fields.ServerSecret,
				RestyClient:  tt.fields.RestyClient,
			}
			if err := j.Get(tt.args.url, tt.args.params, tt.args.resBody); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJiGou_GetPublicParams(t *testing.T) {
	type fields struct {
		AppId        string
		ServerSecret string
		RestyClient  *resty.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JiGou{
				AppId:        tt.fields.AppId,
				ServerSecret: tt.fields.ServerSecret,
				RestyClient:  tt.fields.RestyClient,
			}
			if got := j.GetPublicParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPublicParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJiGou_GetRoomNumbers1(t *testing.T) {
	type fields struct {
		AppId        string
		ServerSecret string
		RestyClient  *resty.Client
	}
	type args struct {
		roomId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RoomNumbers
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JiGou{
				AppId:        tt.fields.AppId,
				ServerSecret: tt.fields.ServerSecret,
				RestyClient:  tt.fields.RestyClient,
			}
			got, err := j.GetRoomNumbers(tt.args.roomId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoomNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoomNumbers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJiGou_GetToken(t *testing.T) {
	type fields struct {
		AppId        string
		ServerSecret string
		RestyClient  *resty.Client
	}
	type args struct {
		userId string
		roomId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JiGou{
				AppId:        tt.fields.AppId,
				ServerSecret: tt.fields.ServerSecret,
				RestyClient:  tt.fields.RestyClient,
			}
			got, err := j.GetToken(tt.args.userId, tt.args.roomId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJiGou_SendCustomCommand(t *testing.T) {
	type fields struct {
		AppId        string
		ServerSecret string
		RestyClient  *resty.Client
	}
	type args struct {
		roomId     string
		fromUserId string
		toUserId   []string
		message    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SendCustomCommand
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JiGou{
				AppId:        tt.fields.AppId,
				ServerSecret: tt.fields.ServerSecret,
				RestyClient:  tt.fields.RestyClient,
			}
			got, err := j.SendCustomCommand(tt.args.roomId, tt.args.fromUserId, tt.args.toUserId, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendCustomCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendCustomCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJiGouClient(t *testing.T) {
	type args struct {
		appid        string
		serverSecret string
	}
	tests := []struct {
		name string
		args args
		want *JiGou
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJiGouClient(tt.args.appid, tt.args.serverSecret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJiGouClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateSignature(t *testing.T) {
	type args struct {
		appId          string
		serverSecret   string
		signatureNonce string
		timeStamp      int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateSignature(tt.args.appId, tt.args.serverSecret, tt.args.signatureNonce, tt.args.timeStamp); got != tt.want {
				t.Errorf("generateSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readJigouConfig(t *testing.T) {
	type args struct {
		v viper.Viper
	}
	tests := []struct {
		name   string
		args   args
		wantJg jigouConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotJg := readJigouConfig(tt.args.v); !reflect.DeepEqual(gotJg, tt.wantJg) {
				t.Errorf("readJigouConfig() = %v, want %v", gotJg, tt.wantJg)
			}
		})
	}
}
