/*
File Name:  redisCmdHandler_test.go.go
Description:
Author:      Chenghu
Date:       2023/8/21 16:28
Change Activity:
*/
package redisDb

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRedisCmdBindName(t *testing.T) {
	type args struct {
		str  string
		args map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "haha", args: args{str: "lrange pinky1:{{name}}   {{value}} ", args: map[string]any{"name": "sssw", "value": []int{2, 3, 4}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RedisCmdBindName(tt.args.str, tt.args.args)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCmdBindName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisCmdBindName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
