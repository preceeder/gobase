/*
File Name:  tm_test.go.go
Description:  默认时区 东8   utc+8   gmt+8
Author:      Chenghu
Date:       2023/9/5 15:59
Change Activity:
*/
package datetimeh

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"reflect"
	"testing"
)

func TestDayNextSecond(t *testing.T) {
	type args struct {
		zone []string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: "", args: args{
			zone: []string{"GMT+4"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DayNextSecond(tt.args.zone...)
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("DayNextSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMouthNextSecond(t *testing.T) {
	type args struct {
		zone []string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MouthNextSexond(tt.args.zone...)
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("MouthNextSexond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekNextSecond(t *testing.T) {
	type args struct {
		zone []string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
		{name: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeekNextSecond(tt.args.zone...)
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("WeekNextSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeZoneHandler(t *testing.T) {
	type args struct {
		zone []string
	}
	tests := []struct {
		name string
		args args
		want carbon.Carbon
	}{
		// TODO: Add test cases.
		{name: "", args: args{zone: []string{"GMT+2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timeZoneHandler(tt.args.zone...)
			got.Now()
			fmt.Println(got.Now().EndOfWeek(), carbon.Now())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeZoneHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
