/*
File Name:  stringH_test.go.go
Description:
Author:      Chenghu
Date:       2023/8/23 16:57
Change Activity:
*/
package utils

import (
	"fmt"
	"testing"
)

func TestGenterWithoutRepetitionStr(t *testing.T) {
	type args struct {
		strl int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "", args: args{strl: 10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenterWithoutRepetitionStr(tt.args.strl)
			if got != tt.want {
				t.Errorf("GenterWithoutRepetitionStr() = %v, want %v", got, tt.want)
			}
			fmt.Println(got)
		})
	}
}

func TestRandStr(t *testing.T) {
	type args struct {
		str_len int
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
			if got := RandStr(tt.args.str_len); got != tt.want {
				t.Errorf("RandStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrBindName(t *testing.T) {
	type args struct {
		str     string
		args    map[string]any
		spacing []byte
	}
	tests := []struct {
		name      string
		args      args
		wantTemPs string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTemPs, err := StrBindName(tt.args.str, tt.args.args, tt.args.spacing)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrBindName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTemPs != tt.wantTemPs {
				t.Errorf("StrBindName() gotTemPs = %v, want %v", gotTemPs, tt.wantTemPs)
			}
		})
	}
}
