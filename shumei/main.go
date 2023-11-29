/*
File Name:  main.go
Description:
Author:      Chenghu
Date:       2023/8/23 11:25
Change Activity:
*/
package shumei

import (
	"github.com/bytedance/sonic"
	mapset "github.com/deckarep/golang-set"
	"github.com/go-resty/resty/v2"
	"github.com/preceeder/gobase/utils"
	"github.com/preceeder/gobase/utils/datetimeh"
	"github.com/spf13/viper"
	"log/slog"
	"net/url"
	"strings"
	"time"
)

var ShumeiClient *ShuMei

var imageLangSet = mapset.NewSet("zh", "en", "ar")
var textLangSet = mapset.NewSet("zh", "en", "ar", "hi", "es", "fr", "ru", "pt", "id", "de", "ja", "tr", "vi", "it", "th", "tl", "ko", "ms", "auto")
var voiceLangSet = mapset.NewSet("zh", "en", "ar", "hi", "es", "fr", "ru", "pt", "id", "de", "ja", "tr", "vi", "it", "th", "tl", "ko", "ms")

type ShumeiUrl struct {
	VideoStreamCloseUrl   string `json:"videoStreamCloseUrl"`   // 视频流检查关闭的数美url
	VoiceStreamCloseUrl   string `json:"voiceStreamCLoseUrl"`   // 语音流检查关闭的数美url
	ImageUrl              string `json:"imageUrl"`              // 同步图片的检查的数美url
	TextUrl               string `json:"textUrl"`               // 文本检查的数美url
	VoiceUrl              string `json:"voiceUrl"`              // 语音文件检查的数美url
	AsyncVoiceUrl         string `json:"asyncVoiceUrl"`         // 异步语音文件检查的数美url
	AsyncVoiceCallBackUrl string `json:"asyncVoiceCallBackUrl"` // 异步语音文件检查回调url
	AsyncVideoUrl         string `json:"asyncVideoUrl"`         // 视频文件检查的数美url
	AsyncVideoCallBackUrl string `json:"asyncVideoCallBackUrl"` // 视频文件检查回调url
	VoiceStreamUrl        string `json:"voiceStreamUrl"`        // 音频流检查url
	VideoStreamUrl        string `json:"videoStreamUrl"`        // 视频流检查url
}
type ShumeiConfig struct {
	AppId          string    `json:"appid"`
	AccessKey      string    `json:"accessKey"`
	CdnUrl         string    `json:"cdnUrl"`         // cdn的url
	TokenPrefix    string    `json:"tokenPrefix"`    // 用户token的前缀
	CallBackDomain string    `json:"callBackDomain"` // 回调的 url
	ShumeiUrl      ShumeiUrl `json:"shumeiUrl"`
}

func initShumei(config ShumeiConfig) {
	client, err := NewShuMei(config.AppId, config.AccessKey,
		OptionWithTokenPrefix(config.TokenPrefix),
		OptionWithCdnUrl(config.CdnUrl),
		OptionWithCallBackDomain(config.CallBackDomain))
	if err != nil {
		slog.Error("数美初始化失败", "error", err.Error())
		panic("数美初始化失败")
	}
	ShumeiClient = client
	ShumeiClient.ShumeiUrl = config.ShumeiUrl
}

// 使用 viper读取的配置初始化
func InitShumeiWithViperConfig(config viper.Viper) {
	shumeisConfig := ShumeiConfig{}
	utils.ReadViperConfig(config, "shumei", &shumeisConfig)
	initShumei(shumeisConfig)
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
	CdnUrl           string    // 资源的 url 最后不加 /
	CallBackDomain   string    // 回调域名  url 最后不加 /
	ShumeiUrl        ShumeiUrl // 数美接口的urls
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

	return sh, nil
}

func OptionWithTokenPrefix(t string) func(*ShuMei) error {
	return func(s *ShuMei) error {
		s.TokenPrefix = t
		return nil
	}
}

func OptionWithCdnUrl(t string) func(mei *ShuMei) error {
	return func(s *ShuMei) error {
		s.CdnUrl = t
		return nil
	}
}

func OptionWithCallBackDomain(t string) func(mei *ShuMei) error {
	return func(s *ShuMei) error {
		s.CallBackDomain = t
		return nil
	}
}
func OptionWithUrl(t ShumeiUrl) func(mei *ShuMei) error {
	return func(s *ShuMei) error {
		s.ShumeiUrl = t
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
		imageUrl, _ = url.JoinPath(s.CdnUrl, imageUrl)
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
func (s ShuMei) AsyncImage(ctx utils.Context, p ShumeiAsyncImage) (bool, *PublicLongResponse) {
	turl := s.ShumeiUrl.ImageUrl //  "http://api-img-sh.fengkongcloud.com/image/v4"

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

	res := &PublicLongResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error(), "requestId", ctx.RequestId)
		return false, nil
	}

	if res.Code != 1100 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "code", res.Code, "requestId", ctx.RequestId)
		return false, nil
	}
	return true, res
}

func (s ShuMei) Image(ctx utils.Context, p ShumeiImage) bool {
	//turl := "http://api-img-xjp.fengkongcloud.com/image/v4"
	turl := s.ShumeiUrl.ImageUrl //"http://api-img-sh.fengkongcloud.com/image/v4"

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
		slog.Error("shumei image request", "error", err.Error(), "requestId", ctx.RequestId)
		return true
	}

	if res.Code == 1100 {
		if res.RiskLevel == "REJECT" {
			return false
		}
	}
	return true
}

func (s ShuMei) Text(ctx utils.Context, p ShumeiText) bool {
	turl := s.ShumeiUrl.TextUrl //"http://api-text-sh.fengkongcloud.com/text/v4"
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
		slog.Error("shumei image request", "error", err.Error(), "requestId", ctx.RequestId)
		return true
	}
	if res.Code == 1100 {
		if res.RiskLevel == "REJECT" {
			return false
		}
	}
	return true
}

// 只支持 url的 同步
func (s ShuMei) VoiceFile(ctx utils.Context, p ShumeiVoiceFile) bool {
	turl := s.ShumeiUrl.VoiceUrl // "http://api-audio-sh.fengkongcloud.com/audiomessage/v4"
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
		slog.Error("shumei image request", "error", err.Error(), "requestId", ctx.RequestId)
		return true
	}
	if res.Code == 1100 {
		if res.Detail.RiskLevel == "REJECT" {
			return false
		}
	}
	return true
}

func (s ShuMei) AsyncVoiceFile(ctx utils.Context, p ShumeiVoiceFile) bool {
	turl := s.ShumeiUrl.AsyncVoiceUrl //"http://api-audio-sh.fengkongcloud.com/audio/v4"
	data := map[string]interface{}{
		"tokenId": s.tokenHandler(p.UserId),
		"lang":    s.voiceLangeHandler(p.Lang),
	}
	if len(p.CallbackParams) > 0 && p.NeedCallback {
		data["passThrough"] = p.CallbackParams
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
	if p.NeedCallback {
		payload["callback"] = s.ShumeiUrl.AsyncVoiceCallBackUrl
	}

	res := &PublicShortResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error(), "requestId", ctx.RequestId)
		return false
	}
	if res.Code != 1100 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "code", res.Code, "requestId", ctx.RequestId)
		return false
	}
	return true
}

func (s ShuMei) AsyncVideoFile(ctx utils.Context, p ShumeiAsyncVideoFile) bool {
	//上海节点
	turl := s.ShumeiUrl.AsyncVideoUrl // "http://api-video-sh.fengkongcloud.com/video/v4"
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
		"callback":  s.ShumeiUrl.AsyncVideoCallBackUrl,
		"data":      data,
	}
	res := &VideoFileResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei image request", "error", err.Error(), "requestId", ctx.RequestId)
		return false
	}
	if res.Code != 1100 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "code", res.Code, "requestId", ctx.RequestId)
		return false
	}
	return true
}

// 回调路径处理
func (s ShuMei) HandlerCallBackUrl(urlStr string) string {
	if !(strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")) {
		urlStr, _ = url.JoinPath(s.CallBackDomain, urlStr)
	}
	return urlStr
}

// 音频流检查
func (s ShuMei) AudioStream(ctx utils.Context, p ShumeiAsyncAudioStream) (bool, *AudioStreamResponse) {
	turl := s.ShumeiUrl.VoiceStreamUrl //"http://api-audiostream-sh.fengkongcloud.com/audiostream/v4"
	data := map[string]interface{}{
		"tokenId":          s.tokenHandler(p.UserId),
		"lang":             s.voiceLangeHandler(p.Lang),
		"btId":             utils.GenterWithoutRepetitionStr(16),
		"streamType":       p.StreamType,
		"returnAllText":    p.ReturnAllText,
		"room":             p.RoomId,
		"returnFinishInfo": 1,
		"audioDetectStep":  p.AudioDetectStep,
		"extra":            map[string]any{"passThrough": p.ThroughParams},
	}

	for k, v := range p.RtcParams {
		data[k] = v
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
		"callback":  s.HandlerCallBackUrl(p.Callback),
		"data":      data,
	}

	res := &AudioStreamResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei AsyncVoiceFile request", "error", err.Error(), "requestId", ctx.RequestId)
		return false, nil
	}
	if res.Code != 1100 {
		slog.Error("AudioStream", "error", res.Message, "requestId", res.RequestID, "code", res.Code, "requestId", ctx.RequestId)
		return false, nil
	} else if res.Code == 1100 && res.Detail.Errorcode != 0 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "errorCode", res.Detail.Errorcode, "requestId", ctx.RequestId)
		return false, nil
	}

	return true, res
}

// 视频流检查
func (s ShuMei) VideoStream(ctx utils.Context, p ShumeiAsyncVideoStream) (bool, *AudioStreamResponse) {
	turl := s.ShumeiUrl.VideoStreamUrl //"http://api-videostream-sh.fengkongcloud.com/videostream/v4"
	data := map[string]interface{}{
		"tokenId":          s.tokenHandler(p.UserId),
		"lang":             s.imageLangHandler(p.Lang),
		"btId":             utils.GenterWithoutRepetitionStr(16),
		"streamType":       p.StreamType,
		"room":             p.RoomId,
		"returnFinishInfo": p.ReturnFinishInfo,
		"detectFrequency":  p.DetectFrequency, // 通知的频次   秒/次
		//"audioDetectStep":  20,
		"extra": map[string]any{"passThrough": p.ThroughParams},
	}

	// rtc 参数
	for k, v := range p.RtcParams {
		data[k] = v
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
		"accessKey":     s.AccessKey,
		"appId":         s.AppId,
		"eventId":       p.EventId,
		"imgType":       p.VideoType,
		"audioType":     p.VoiceType,
		"imgCallback":   s.HandlerCallBackUrl(p.ImgCallback),
		"audioCallback": s.HandlerCallBackUrl(p.AudioCallback),
		"data":          data,
	}

	res := &AudioStreamResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei AsyncVoiceFile request", "error", err.Error(), "requestId", ctx.RequestId)
		return false, nil
	}
	if res.Code != 1100 {
		slog.Error("AudioStream", "error", res.Message, "requestId", res.RequestID, "code", res.Code, "requestId", ctx.RequestId)
		return false, nil
	} else if res.Code == 1100 && res.Detail.Errorcode != 0 {
		slog.Error("AsyncVoiceFile", "error", res.Message, "requestId", res.RequestID, "errorCode", res.Detail.Errorcode, "requestId", ctx.RequestId)
		return false, nil
	}

	return true, res
}

/** 流关闭接口
 * @param requestId string 请求id
 * @ltype string 类型 voice｜video
 */
func (s ShuMei) CloseStreamCheck(ctx utils.Context, requestId string, ltype string) (bool, *CloseStreamResponse) {
	turl := ""
	if ltype == "video" {
		turl = s.ShumeiUrl.VideoStreamCloseUrl
	} else if ltype == "voice" {
		turl = s.ShumeiUrl.VoiceStreamCloseUrl
	}

	payload := map[string]interface{}{
		"accessKey": s.AccessKey,
		"requestId": requestId,
	}

	res := &CloseStreamResponse{}
	_, err := s.Send(turl, payload, res)
	if err != nil {
		slog.Error("shumei close stream faild", "error", err.Error(), "requestId", ctx.RequestId)
		return false, nil
	}
	return true, res

}

func (s ShuMei) TianWang(ctx utils.Context, p TianWangParams) map[string]any {
	data := map[string]any{
		"accessKey": s.AccessKey,
		"appId":     s.AppId,
		"eventId":   p.EventId,
		"data": map[string]any{
			"tokenId":    p.TokenId,
			"ip":         p.Ip,
			"timestamp":  datetimeh.Now().TimestampMilli(),
			"deviceId":   p.SmDeviceId,
			"phone":      p.Phone,
			"os":         p.Channel,
			"appVersion": p.Version,
			"type":       p.RegisterMethod,
		},
	}
	response, err := s.Send("http://api-skynet-bj.fengkongcloud.com/v4/event", data, nil)
	if err != nil {
		slog.Error("tianwang 事件接口访问失败", "error", err.Error(), "requestId", ctx.RequestId)
	}
	return response
}
