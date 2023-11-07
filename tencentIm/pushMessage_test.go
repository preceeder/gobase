//   File Name:  pushMessage_test.go.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/6 16:57
//    Change Activity:

package tencentIm

import (
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
		{name: "", args: args{ctx: utils.Context{}, serverName: "SendMsg", fromId: "u2_42", toId: "u2_1162203", content: TextContent{Text: "haha"}}},
	}

	InitWithStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//res := &CommonResponse{}
			if err := SendImMessage(tt.args.ctx, tt.args.serverName, tt.args.fromId, tt.args.toId, tt.args.content, tt.args.cloudCustomData, tt.args.sendMsgControl, tt.args.forbidCallbackControl, tt.args.syncOtherMachine, nil, nil); (err != nil) != tt.wantErr {
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
		{name: "", args: args{ctx: utils.Context{}, serverName: "BatchSendMsg", fromId: "u2_42", toId: []string{"u2_1162203"}, content: TextContent{Text: "haha"}}},
	}

	InitWithStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := BatchCommonResponse{}
			if err := SendBatchImMessage(tt.args.ctx, tt.args.serverName, tt.args.fromId, tt.args.toId, tt.args.content, tt.args.cloudCustomData, tt.args.sendMsgControl, tt.args.syncOtherMachine, nil, &res); (err != nil) != tt.wantErr {
				t.Errorf("SendImMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
