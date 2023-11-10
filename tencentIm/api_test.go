//   File Name:  pushMessage_test.go.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/6 16:57
//    Change Activity:

package tencentIm

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/golang-module/carbon/v2"
	"github.com/preceeder/gobase/utils"
	"testing"
)

func TestSendImMessage(t *testing.T) {
	type args struct {
		ctx                   utils.Context
		serverName            string
		fromId                string
		toId                  string
		content               MsgContent
		cloudCustomData       string
		sendMsgControl        []string
		forbidCallbackControl []string
		syncOtherMachine      int
		offLineData           OfflinePushInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "", args: args{ctx: utils.Context{}, fromId: "u2_42", toId: "u2_1162203", content: TextContent{Text: "haha"}}},
	}

	InitWithStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//res := &CommonResponse{}
			if err := SendImMessage(tt.args.ctx, tt.args.fromId, tt.args.toId, tt.args.content, tt.args.cloudCustomData, tt.args.sendMsgControl, tt.args.forbidCallbackControl, tt.args.syncOtherMachine, nil, nil); (err != nil) != tt.wantErr {
				t.Errorf("SendImMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendBatchImMessage(t *testing.T) {
	type args struct {
		ctx              utils.Context
		serverName       string
		fromId           string
		toId             []string
		content          MsgContent
		cloudCustomData  string
		sendMsgControl   []string
		syncOtherMachine int
		offLineData      OfflinePushInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "", args: args{ctx: utils.Context{}, fromId: "u2_42", toId: []string{"u2_1162203"}, content: TextContent{Text: "haha"}}},
	}

	InitWithStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := BatchCommonResponse{}
			if err := SendBatchImMessage(tt.args.ctx, tt.args.fromId, tt.args.toId, tt.args.content, tt.args.cloudCustomData, tt.args.sendMsgControl, tt.args.syncOtherMachine, nil, &res); (err != nil) != tt.wantErr {
				t.Errorf("SendImMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteRecentContact(t *testing.T) {
	type args struct {
		ctx         utils.Context
		fromAccount string
		toAccount   string
		toGroupid   string
		htype       int
		clearRamble int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteRecentContact(tt.args.ctx, tt.args.fromAccount, tt.args.toAccount, tt.args.toGroupid, tt.args.htype, tt.args.clearRamble); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRecentContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMsgWithdraw(t *testing.T) {
	type args struct {
		ctx         utils.Context
		fromAccount string
		toAccount   string
		msgKey      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MsgWithdraw(tt.args.ctx, tt.args.fromAccount, tt.args.toAccount, tt.args.msgKey); (err != nil) != tt.wantErr {
				t.Errorf("MsgWithdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryHistoryMessage(t *testing.T) {
	type args struct {
		ctx             utils.Context
		operatorAccount string
		peerAccount     string
		lastMsgKey      string
		maxCnt          int
		minTime         int64
		maxTime         int64
	}
	now := carbon.Now()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "", args: args{ctx: utils.Context{}, operatorAccount: "u2_1162203", peerAccount: "u2_42", maxCnt: 100, minTime: now.SubDays(10).Timestamp(), maxTime: now.Timestamp()}},
	}
	TencentImConfig.ImHost = "https://console.tim.qq.com"
	TencentImConfig.UseSha = "ECDSA-SHA256"
	TencentImConfig.Expire = 3600 * 24 * 360

	InitWithStruct()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd, err := QueryHistoryMessage(tt.args.ctx, tt.args.operatorAccount, tt.args.peerAccount, tt.args.lastMsgKey, tt.args.maxCnt, tt.args.minTime, tt.args.maxTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryHistoryMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			dert, _ := sonic.MarshalString(dd)
			fmt.Println(dert)
		})
	}
}
