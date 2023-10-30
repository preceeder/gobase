/*
File Name:  logs_test.go.go
Description:
Author:      Chenghu
Date:       2023/9/3 11:16
Change Activity:
*/
package logs

import (
	"log/slog"
	"testing"
)

func TestInitLogWithStruct(t *testing.T) {
	type args struct {
		cfg LogConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "ss", args: args{cfg: LogConfig{
			DebugFileName: "logs/log.txt",
			InfoFileName:  "logs/log.txt",
			WarnFileName:  "logs/log.txt",
			MaxSize:       10,
			MaxAge:        1,
			MaxBackups:    3,
			StdOut:        "1",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitLogWithStruct(tt.args.cfg)
			slog.Info("sdee", "nid", "sddddd")
		})
	}
}
