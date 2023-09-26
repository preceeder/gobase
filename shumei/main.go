/*
File Name:  main.py
Description:
Author:      Chenghu
Date:       2023/8/23 11:25
Change Activity:
*/
package shumei

import (
	"fmt"
	"github.com/bytedance/sonic"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-resty/resty/v2"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"
)

var ShumeiClient *ShuMei

var imageLangSet = mapset.NewSet("zh", "en", "ar")
var textLangSet = mapset.NewSet("zh", "en", "ar", "hi", "es", "fr", "ru", "pt", "id", "de", "ja", "tr", "vi", "it", "th", "tl", "ko", "ms", "auto")
var voiceLangSet = mapset.NewSet("zh", "en", "ar", "hi", "es", "fr", "ru", "pt", "id", "de", "ja", "tr", "vi", "it", "th", "tl", "ko", "ms")

type ShumeiConfig struct {
	AppId       string `json:"appid"`
	accessKey   string `json:"accessKey"`
	CdnUrl      string `json:"cdnUrl"`      // cdn的url
	TokenPrefix string `json:"tokenPrefix"` // 用户token的前缀
}

func initShumei(config ShumeiConfig) {
	client, err := NewShuMei(config.AppId, config.accessKey, OptionWithTokenPrefix(config.TokenPrefix), OptionWithBaseUrl(config.CdnUrl))
	if err != nil {
		slog.Error("数美初始化失败", "error", err.Error())
		panic("数美初始化失败")
	}
	ShumeiClient = client
}

// 使用 viper读取的配置初始化
func InitShumeiWithViperConfig(config viper.Viper) {
	shumeisConfig := readRedisConfig(config)
	initShumei(shumeisConfig)
}

func readRedisConfig(v viper.Viper) (rs ShumeiConfig) {
	shumei := v.Sub("shumei")
	if shumei == nil {
		fmt.Printf("shumei config is nil")
		os.Exit(1)
	}
	rs = ShumeiConfig{}
	err := shumei.Unmarshal(&rs)
	if err != nil {
		fmt.Printf("shumei config read error: " + err.Error())
		os.Exit(1)
	}
	return
}

type ShuMei struct {
	AppId            string
	AccessKey        string
	DefaultImageType string // 默认值 POLITICS_PORN_AD
	DefaultTextType  string // 默认值 AD
	DefaultVoiceType string // 默认值 PORN_MOAN_AD
	DefaultVideoType string // 默认值 POLITY_EROTIC_ADVERT
	TokenPrefix      string // 用户id的统一前缀
	HttpClient       *resty.Client
	BaseUrl          string
	StreamType       string // 目前默认值 ZEGO
}

func NewShuMei(appId string, accessKey string, optionals ...func(*ShuMei) error) (*ShuMei, error) {
	tp := resty.New().SetTimeout(time.Duration(5 * time.Second))
	sh := &ShuMei{
		AppId:      appId,
		AccessKey:  accessKey,
		HttpClient: tp,
	}
	for _, op := range optionals {
		err := op(sh)
		if err != nil {
			return nil, err
		}
	}
	if sh.DefaultImageType == "" {
		sh.DefaultImageType = "POLITICS_PORN_AD"
	}
	if sh.DefaultVoiceType == "" {
		sh.DefaultVoiceType = "PORN_MOAN_AD"
	}

	if sh.DefaultTextType == "" {
		sh.DefaultTextType = "AD"
	}
	if sh.DefaultVideoType == "" {
		sh.DefaultVideoType = "POLITY_EROTIC_ADVERT"
	}
	if sh.StreamType == "" {
		sh.StreamType = "ZEGO"
	}

	return sh, nil
}

func OptionWithTokenPrefix(t string) func(*ShuMei) error {
	return func(s *ShuMei) error {
		s.TokenPrefix = t
		return nil
	}
}

func OptionWithBaseUrl(t string) func(mei *ShuMei) error {
	return func(s *ShuMei) error {
		s.BaseUrl = t
		return nil
	}
}

func OptionWithImageType(t string) func(*ShuMei) error {
	return func(s *ShuMei) error {
		s.DefaultImageType = t
		return nil
	}
}

func OptionWithTextType(t string) func(*ShuMei) error {
	return func(s *ShuMei) error {
		s.DefaultTextType = t
		return nil
	}
}

func (s ShuMei) Send(url string, body any, response any) (map[string]any, error) {
	res, err := s.HttpClient.R().
		SetResult(response).
		ForceContentType("application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		EnableTrace().
		SetBody(body).
		Post(url)
	fmt.Println(res.Result())
	if response == nil {
		var data map[string]interface{}
		_ = sonic.Unmarshal(res.Body(), &data)
		return data, err
	} else {
		_ = sonic.Unmarshal(res.Body(), response)
		return nil, err
	}
}

func (s ShuMei) urlHandler(imageUrl string) string {
	if !(strings.HasPrefix(imageUrl, "http://") || strings.HasPrefix(imageUrl, "https://")) {
		imageUrl, _ = url.JoinPath(s.BaseUrl, imageUrl)
	}
	return imageUrl
}

func (s ShuMei) tokenHandler(userId string) string {
	if userId == "" {
		return ""
	}
	if !strings.HasPrefix(userId, s.TokenPrefix) {
		userId = s.TokenPrefix + userId
	}
	return userId
}

func (s ShuMei) voiceLangeHandler(lang string) string {
	if !voiceLangSet.Contains(lang) {
		lang = "zh"
	}
	return lang
}

func (s ShuMei) imageLangHandler(lang string) string {
	if !imageLangSet.Contains(lang) {
		lang = "zh"
	}
	return lang
}

func (s ShuMei) textLangHandler(lang string) string {
	if !textLangSet.Contains(lang) {
		lang = "auto"
	}
	return lang
}

// 同步 图片检查
func (s ShuMei) AsyncImage(p ShumeiAsyncImage) bool {
	//turl := "http://api-img-xjp.fengkongcloud.com/image/v4"
	turl := "http://api-img-sh.fengkongcloud.com/image/v4"

	data := map[string]interface{}{
		"img":            s.urlHandler(p.ImageUrl),
		"tokenId":        s.tokenHandler(p.UserId),
		"receiveTokenId": s.tokenHandler(p.ReceiveTokenId),
		"lang":           s.imageLangHandler(p.Lang),
	}
	if p.Ip != "" {
		data["ip"] = p.Ip
	}

	if p.MType == "" {
		p.MType = s.DefaultImageType
	}

	data["extra"] = map[string]any{
		"passThrough": p.ThroughParams,
	}

	payload := map[string]interface{}{
		"accessKey":    s.AccessKey,
		"type":         p.MType,
		"eventId":      "IMAGE",
		"businessType": "FACE",
		"appId":        s.AppId,
		"data":         data,
		"callback":     p.CallBaskUrl,
	}

	res := &PublicShortResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error())
		return false
	}

	if res.Code != 1100 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "code", res.Code)
		return false
	}
	return true
}

func (s ShuMei) Image(p ShumeiImage) bool {
	//turl := "http://api-img-xjp.fengkongcloud.com/image/v4"
	turl := "http://api-img-sh.fengkongcloud.com/image/v4"

	data := map[string]interface{}{
		"img":            s.urlHandler(p.ImageUrl),
		"tokenId":        s.tokenHandler(p.UserId),
		"receiveTokenId": s.tokenHandler(p.ReceiveTokenId),
		"lang":           s.imageLangHandler(p.Lang),
	}
	if p.Ip != "" {
		data["ip"] = p.Ip
	}

	if p.MType == "" {
		p.MType = s.DefaultImageType
	}

	payload := map[string]interface{}{
		"accessKey":    s.AccessKey,
		"type":         p.MType,
		"eventId":      "IMAGE",
		"businessType": "FACE",
		"appId":        s.AppId,
		"data":         data,
	}

	res := &PublicLongResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error())
		return true
	}

	if res.Code == 1100 {
		if res.RiskLevel == "REJECT" {
			return false
		}
	}
	return true
}

func (s ShuMei) Text(p ShumeiText) bool {
	turl := "http://api-text-sh.fengkongcloud.com/text/v4"
	data := map[string]interface{}{
		"text":    p.Text,
		"tokenId": s.tokenHandler(p.UserId),
		"lang":    s.textLangHandler(p.Lang),
	}
	if p.EventId == "" {
		p.EventId = "text"
	}
	if p.EventId == "message" {
		data["extra"] = map[string]any{"receiveTokenId": s.tokenHandler(p.ReceiveTokenId)}
	}

	if p.MType == "" {
		p.MType = s.DefaultTextType
	}

	payload := map[string]interface{}{
		"accessKey": s.AccessKey,
		"appId":     s.AppId,
		"eventId":   p.EventId,
		"type":      p.MType,
		"data":      data,
	}
	res := &PublicLongResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error())
		return true
	}
	if res.Code == 1100 {
		if res.RiskLevel == "REJECT" {
			return false
		}
	}
	fmt.Println(res)
	return true
}

// 只支持 url的 同步
func (s ShuMei) VoiceFile(p ShumeiVoiceFile) bool {
	turl := "http://api-audio-sh.fengkongcloud.com/audiomessage/v4"
	data := map[string]interface{}{
		"tokenId": s.tokenHandler(p.UserId),
	}
	if p.EventId == "" {
		p.EventId = "default"
	}
	if p.MType == "" {
		p.MType = s.DefaultVoiceType
	}

	payload := map[string]interface{}{
		"accessKey":   s.AccessKey,
		"appId":       s.AppId,
		"eventId":     p.EventId,
		"type":        p.MType,
		"contentType": "URL",
		"content":     p.VoiceUrl,
		"data":        data,
		"btId":        utils.GenterWithoutRepetitionStr(16),
	}

	res := &VoiceFileResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error())
		return true
	}
	if res.Code == 1100 {
		if res.Detail.RiskLevel == "REJECT" {
			return false
		}
	}
	return true
}

func (s ShuMei) AsyncVoiceFile(p ShumeiVoiceFile) bool {
	//上海
	turl := "http://api-audio-sh.fengkongcloud.com/audio/v4"
	data := map[string]interface{}{
		"tokenId": s.tokenHandler(p.UserId),
		"lang":    s.voiceLangeHandler(p.Lang),
	}

	if p.EventId == "" {
		p.EventId = "default"
	}
	if p.MType == "" {
		p.MType = s.DefaultVoiceType
	}

	payload := map[string]interface{}{
		"accessKey":   s.AccessKey,
		"appId":       s.AppId,
		"eventId":     p.EventId,
		"type":        p.MType,
		"contentType": "URL",
		"content":     p.VoiceUrl,
		"data":        data,
		"btId":        utils.GenterWithoutRepetitionStr(16),
	}

	res := &PublicShortResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error())
		return false
	}
	if res.Code != 1100 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "code", res.Code)
		return false
	}
	return true
}

func (s ShuMei) AsyncVideoFile(p ShumeiAsyncVideoFile) bool {
	//上海节点
	turl := "http://api-video-sh.fengkongcloud.com/video/v4"
	data := map[string]interface{}{
		"tokenId": s.tokenHandler(p.UserId),
		"lang":    s.imageLangHandler(p.Lang),
		"btId":    utils.GenterWithoutRepetitionStr(16),
		"url":     s.urlHandler(p.VideoUrl),
		"extra":   map[string]any{"passThrough": p.ThroughParams},
	}

	if p.EventId == "" {
		p.EventId = "default"
	}
	if p.VideoType == "" {
		p.VideoType = s.DefaultVideoType
	}

	if p.VoiceType == "" {
		p.VoiceType = s.DefaultVoiceType
	}

	payload := map[string]interface{}{
		"accessKey": s.AccessKey,
		"appId":     s.AppId,
		"eventId":   p.EventId,
		"imgType":   p.VideoType,
		"audioType": p.VoiceType,
		"callback":  p.Callback,
		"data":      data,
	}
	res := &VideoFileResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error())
		return false
	}
	if res.Code != 1100 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "code", res.Code)
		return false
	}
	return true
}

// 音频流检查
func (s ShuMei) AudioStream(p ShumeiAsyncAudioStream) (bool, *AudioStreamResponse) {
	turl := "http://api-audiostream-sh.fengkongcloud.com/audiostream/v4"
	data := map[string]interface{}{
		"tokenId":    s.tokenHandler(p.UserId),
		"lang":       s.voiceLangeHandler(p.Lang),
		"btId":       utils.GenterWithoutRepetitionStr(16),
		"streamType": s.StreamType,
		"zegoParam": map[string]any{
			"tokenId":  p.TokenId,
			"streamId": p.StreamId,
			"roomId":   p.RoomId,
			"testEnv":  p.TestEnv, /// 使用正式环境
		},
		"returnAllText":    p.ReturnAllText,
		"room":             p.RoomId,
		"returnFinishInfo": 1,
		"audioDetectStep":  p.AudioDetectStep,
		"extra":            map[string]any{"passThrough": p.ThroughParams},
	}

	if p.EventId == "" {
		p.EventId = "default"
	}

	if p.VoiceType == "" {
		p.VoiceType = s.DefaultVoiceType
	}

	payload := map[string]interface{}{
		"accessKey": s.AccessKey,
		"appId":     s.AppId,
		"eventId":   p.EventId,
		"type":      p.VoiceType,
		"callback":  p.Callback,
		"data":      data,
	}

	res := &AudioStreamResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei AsyncVoiceFile request", "error", err.Error())
		return false, nil
	}
	if res.Code != 1100 {
		slog.Error("AudioStream", "error", res.Message, "requestId", res.RequestID, "code", res.Code)
		return false, nil
	} else if res.Code == 1100 && res.Detail.Errorcode != 0 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "errorCode", res.Detail.Errorcode)
		return false, nil
	}

	return true, res
}

// 视频流检查
func (s ShuMei) VideoStream(p ShumeiAsyncVideoStream) (bool, *AudioStreamResponse) {
	turl := "http://api-videostream-sh.fengkongcloud.com/videostream/v4"
	data := map[string]interface{}{
		"tokenId":    s.tokenHandler(p.UserId),
		"lang":       s.imageLangHandler(p.Lang),
		"btId":       utils.GenterWithoutRepetitionStr(16),
		"streamType": s.StreamType,
		"zegoParam": map[string]any{
			"tokenId":  p.TokenId,
			"streamId": p.StreamId,
			"roomId":   p.RoomId,
			"testEnv":  p.TestEnv, /// 使用正式环境
		},
		"room":             p.RoomId,
		"returnFinishInfo": p.ReturnFinishInfo,
		"detectFrequency":  p.DetectFrequency, // 通知的频次   秒/次
		//"audioDetectStep":  20,
		"extra": map[string]any{"passThrough": p.ThroughParams},
	}

	if p.EventId == "" {
		p.EventId = "default"
	}

	if p.VideoType == "" {
		p.VideoType = s.DefaultVideoType
	}

	if p.VoiceType == "" {
		p.VoiceType = s.DefaultVoiceType
	}

	payload := map[string]interface{}{
		"accessKey":   s.AccessKey,
		"appId":       s.AppId,
		"eventId":     p.EventId,
		"imgType":     p.VideoType,
		"audioType":   p.VoiceType,
		"imgCallback": p.ImgCallback,
		"data":        data,
	}

	res := &AudioStreamResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei AsyncVoiceFile request", "error", err.Error())
		return false, nil
	}
	if res.Code != 1100 {
		slog.Error("AudioStream", "error", res.Message, "requestId", res.RequestID, "code", res.Code)
		return false, nil
	} else if res.Code == 1100 && res.Detail.Errorcode != 0 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "errorCode", res.Detail.Errorcode)
		return false, nil
	}

	return true, res
}
