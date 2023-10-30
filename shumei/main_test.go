/*
File Name:  main_test.go.go
Description:
Author:      Chenghu
Date:       2023/8/23 14:02
Change Activity:
*/
package shumei

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/preceeder/gobase/utils"
	"testing"
)

func TestShuMei_AsyncImage1(t *testing.T) {
	type fields struct {
		AppId            string
		AccessKey        string
		DefaultImageType string
		DefaultTextType  string
		DefaultVoiceType string
		DefaultVideoType string
		TokenPrefix      string
		HttpClient       *resty.Client
		BaseUrl          string
		StreamType       string
	}
	type args struct {
		p ShumeiAsyncImage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ShuMei{
				AppId:            tt.fields.AppId,
				AccessKey:        tt.fields.AccessKey,
				DefaultImageType: tt.fields.DefaultImageType,
				DefaultTextType:  tt.fields.DefaultTextType,
				DefaultVoiceType: tt.fields.DefaultVoiceType,
				DefaultVideoType: tt.fields.DefaultVideoType,
				TokenPrefix:      tt.fields.TokenPrefix,
				HttpClient:       tt.fields.HttpClient,
				CdnUrl:           tt.fields.BaseUrl,
				StreamType:       tt.fields.StreamType,
			}
			if got := s.AsyncImage(utils.Context{}, tt.args.p); got != tt.want {
				t.Errorf("AsyncImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuMei_AsyncVideoFile(t *testing.T) {
	type fields struct {
		AppId            string
		AccessKey        string
		DefaultImageType string
		DefaultTextType  string
		DefaultVoiceType string
		DefaultVideoType string
		TokenPrefix      string
		HttpClient       *resty.Client
		BaseUrl          string
		StreamType       string
	}
	type args struct {
		p ShumeiAsyncVideoFile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ShuMei{
				AppId:            tt.fields.AppId,
				AccessKey:        tt.fields.AccessKey,
				DefaultImageType: tt.fields.DefaultImageType,
				DefaultTextType:  tt.fields.DefaultTextType,
				DefaultVoiceType: tt.fields.DefaultVoiceType,
				DefaultVideoType: tt.fields.DefaultVideoType,
				TokenPrefix:      tt.fields.TokenPrefix,
				HttpClient:       tt.fields.HttpClient,
				CdnUrl:           tt.fields.BaseUrl,
				StreamType:       tt.fields.StreamType,
			}
			if got := s.AsyncVideoFile(utils.Context{}, tt.args.p); got != tt.want {
				t.Errorf("AsyncVideoFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuMei_AsyncVoiceFile(t *testing.T) {
	type fields struct {
		AppId            string
		AccessKey        string
		DefaultImageType string
		DefaultTextType  string
		DefaultVoiceType string
		DefaultVideoType string
		TokenPrefix      string
		HttpClient       *resty.Client
		BaseUrl          string
		StreamType       string
	}
	type args struct {
		p ShumeiVoiceFile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ShuMei{
				AppId:            tt.fields.AppId,
				AccessKey:        tt.fields.AccessKey,
				DefaultImageType: tt.fields.DefaultImageType,
				DefaultTextType:  tt.fields.DefaultTextType,
				DefaultVoiceType: tt.fields.DefaultVoiceType,
				DefaultVideoType: tt.fields.DefaultVideoType,
				TokenPrefix:      tt.fields.TokenPrefix,
				HttpClient:       tt.fields.HttpClient,
				CdnUrl:           tt.fields.BaseUrl,
				StreamType:       tt.fields.StreamType,
			}
			if got := s.AsyncVoiceFile(utils.Context{}, tt.args.p); got != tt.want {
				t.Errorf("AsyncVoiceFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuMei_AudioStream(t *testing.T) {
	type fields struct {
		AppId            string
		AccessKey        string
		DefaultImageType string
		DefaultTextType  string
		DefaultVoiceType string
		DefaultVideoType string
		TokenPrefix      string
		HttpClient       *resty.Client
		BaseUrl          string
		StreamType       string
	}
	type args struct {
		p ShumeiAsyncAudioStream
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ShuMei{
				AppId:            tt.fields.AppId,
				AccessKey:        tt.fields.AccessKey,
				DefaultImageType: tt.fields.DefaultImageType,
				DefaultTextType:  tt.fields.DefaultTextType,
				DefaultVoiceType: tt.fields.DefaultVoiceType,
				DefaultVideoType: tt.fields.DefaultVideoType,
				TokenPrefix:      tt.fields.TokenPrefix,
				HttpClient:       tt.fields.HttpClient,
				CdnUrl:           tt.fields.BaseUrl,
				StreamType:       tt.fields.StreamType,
			}
			if got, _ := s.AudioStream(utils.Context{}, tt.args.p); got != tt.want {
				t.Errorf("AudioStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuMei_Image(t *testing.T) {

	type args struct {
		p ShumeiImage
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "ss", args: args{p: ShumeiImage{ImageUrl: "", UserId: ""}}},
	}
	s, _ := NewShuMei("", "", OptionWithTokenPrefix("test_"), OptionWithCdnUrl(""))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Image(utils.Context{}, tt.args.p); got != tt.want {
				t.Errorf("Image() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuMei_Text1(t *testing.T) {

	type args struct {
		p ShumeiText
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "ss", args: args{p: ShumeiText{Text: "不是吧", UserId: "2222"}}},
	}

	initShumei(ShumeiConfig{
		AppId:       "",
		AccessKey:   "",
		CdnUrl:      "https://",
		TokenPrefix: "test_",
		ShumeiUrl: ShumeiUrl{
			TextUrl: "http://api-text-xjp.fengkongcloud.com/text/v4",
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShumeiClient.Text(utils.Context{}, tt.args.p)
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuMei_VideoStream(t *testing.T) {
	type fields struct {
		AppId            string
		AccessKey        string
		DefaultImageType string
		DefaultTextType  string
		DefaultVoiceType string
		DefaultVideoType string
		TokenPrefix      string
		HttpClient       *resty.Client
		BaseUrl          string
		StreamType       string
	}
	type args struct {
		p ShumeiAsyncVideoStream
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ShuMei{
				AppId:            tt.fields.AppId,
				AccessKey:        tt.fields.AccessKey,
				DefaultImageType: tt.fields.DefaultImageType,
				DefaultTextType:  tt.fields.DefaultTextType,
				DefaultVoiceType: tt.fields.DefaultVoiceType,
				DefaultVideoType: tt.fields.DefaultVideoType,
				TokenPrefix:      tt.fields.TokenPrefix,
				HttpClient:       tt.fields.HttpClient,
				CdnUrl:           tt.fields.BaseUrl,
				StreamType:       tt.fields.StreamType,
			}
			if got, _ := s.VideoStream(utils.Context{}, tt.args.p); got != tt.want {
				t.Errorf("VideoStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuMei_VoiceFile(t *testing.T) {
	type fields struct {
		AppId            string
		AccessKey        string
		DefaultImageType string
		DefaultTextType  string
		DefaultVoiceType string
		DefaultVideoType string
		TokenPrefix      string
		HttpClient       *resty.Client
		BaseUrl          string
		StreamType       string
	}
	type args struct {
		p ShumeiVoiceFile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ShuMei{
				AppId:            tt.fields.AppId,
				AccessKey:        tt.fields.AccessKey,
				DefaultImageType: tt.fields.DefaultImageType,
				DefaultTextType:  tt.fields.DefaultTextType,
				DefaultVoiceType: tt.fields.DefaultVoiceType,
				DefaultVideoType: tt.fields.DefaultVideoType,
				TokenPrefix:      tt.fields.TokenPrefix,
				HttpClient:       tt.fields.HttpClient,
				CdnUrl:           tt.fields.BaseUrl,
				StreamType:       tt.fields.StreamType,
			}
			if got := s.VoiceFile(utils.Context{}, tt.args.p); got != tt.want {
				t.Errorf("VoiceFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
