/*
File Name:  token_test.go.py
Description:
Author:      Chenghu
Date:       2023/8/21 14:42
Change Activity:
*/
package gobase

import (
	"fmt"
	"reflect"
	"testing"
)

var token = ""

func TestTokenGenerateUsingHs256(t *testing.T) {

	//type TokenInfoData struct {
	//	UserId   string `json:"user_id"`
	//	UserName string `json:"user_name"`
	//}

	type args struct {
		Clam map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "gen", args: args{
			Clam: map[string]string{"name": "nis"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TokenGenerateUsingHs256(tt.args.Clam)
			fmt.Printf("%#v\n", got)
			token = got
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenGenerateUsingHs256() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TokenGenerateUsingHs256() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenParseHs256(t *testing.T) {
	type args struct {
		tokenSecrete string
	}
	tests := []struct {
		name    string
		args    args
		want    *CustomClaims
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "pares", args: args{tokenSecrete: token}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TokenParseHs256(tt.args.tokenSecrete)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenParseHs256() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenParseHs256() got = %v, want %v", got, tt.want)
			}
		})
	}
}
