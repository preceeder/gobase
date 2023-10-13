//   File Name:  binlog_test.go.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/13 16:32
//    Change Activity:

package mysqlDb

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		config BinlogConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "", args: args{
			config: BinlogConfig{
				Addr:     "127.0.0.1:13306",
				Password: "xxxxx",
				User:     "xxxxx",
				Db:       "xxxxx",
			},
		}},
	}
	type Logssd struct {
		ID     int    `json:"id"`
		UserID string `json:"user_id"`
	}

	SetBinlogTable("t_user", []string{"insert", "update", "delete"}, Logssd{}, func(action string, a any, a2 ...any) {
		fmt.Println(action, a)
		fmt.Println(action, a2)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BinLogRun(tt.args.config)
		})
	}
}
