/*
File Name:  allStruct.go
Description:
Author:      Chenghu
Date:       2023/8/23 16:29
Change Activity:
*/
package shumei

// ---------------   请求参数 --------------
type ShumeiImage struct {
	ImageUrl       string
	UserId         string
	ReceiveTokenId string
	MType          string
	Lang           string
	Ip             string
	ThroughParams  map[string]any
	NeedCallBack   bool
	CallBaskUrl    string
}

// 同步
type ShumeiMultiImage struct {
	ImageUrl       []string
	UserId         string
	ReceiveTokenId string
	MType          string
	Lang           string
	Ip             string
	ThroughParams  map[string]any
	NeedCallBack   bool
	CallBaskUrl    string
}

type ShumeiText struct {
	Text           string
	UserId         string
	ReceiveTokenId string
	MType          string
	Lang           string
	EventId        string
}

type ShumeiVoiceFile struct {
	VoiceUrl       string
	UserId         string
	ReceiveTokenId string
	MType          string
	EventId        string
	CallbackUrl    string         // 异步回调需要
	Lang           string         // 异步回调需要
	CallbackParams map[string]any // 异步回调需要
}

// 只有异步
type ShumeiAsyncVideoFile struct {
	VideoUrl       string
	UserId         string
	ReceiveTokenId string
	VideoType      string
	VoiceType      string
	EventId        string
	Lang           string
	CallBackUrl    string
	ThroughParams  map[string]any
}

// 只有异步
type ShumeiAsyncAudioStream struct {
	RtcParams        map[string]any
	StreamType       string // 目前默认值 ZEGO
	UserId           string
	ReceiveTokenId   string
	VoiceType        string
	EventId          string
	Callback         string
	Lang             string
	AudioDetectStep  int
	RoomId           string
	ReturnAllText    int // 0：返回风险等级为非pass的音频片段  1：返回所有风险等级的音频片段   默认0
	ReturnFinishInfo int
	ThroughParams    map[string]any
}

// 只有异步
type ShumeiAsyncVideoStream struct {
	UserId         string
	ReceiveTokenId string
	VideoType      string
	VoiceType      string
	EventId        string
	ImgCallback    string // 视频流只检查 画面
	AudioCallback  string // 音频画面
	ReturnAllImg   int
	ReturnAllText  int
	//Callback       string
	ReturnFinishInfo int
	Lang             string
	RtcParams        map[string]any
	StreamType       string // 目前默认值 ZEGO
	RoomId           string
	DetectFrequency  int // 检测频次
	ThroughParams    map[string]any
}

// ---------------- 响应数据 --------------

// 共长返回数据
type PublicLongResponse struct {
	RequestID       string           `json:"requestId"`
	Code            int              `json:"code"`
	Message         string           `json:"message"`
	RiskLevel       string           `json:"riskLevel"`
	RiskLabel1      string           `json:"riskLabel1"`
	RiskLabel2      string           `json:"riskLabel2"`
	RiskLabel3      string           `json:"riskLabel3"`
	RiskDescription string           `json:"riskDescription"`
	RiskDetail      map[string]any   `json:"riskDetail"`
	AuxInfo         map[string]any   `json:"auxInfo"`
	AllLabels       []map[string]any `json:"allLabels"`
	BusinessLabels  []map[string]any `json:"businessLabels"`
	TokenLabels     map[string]any   `json:"tokenLabels"`
}

// 公共短返回数据
type PublicShortResponse struct {
	RequestID string `json:"requestId"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
}

type VideoFileResponse struct {
	RequestID string `json:"requestId"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	BtId      string `json:"btId"`
}

type AudioStreamResponse struct {
	RequestID string                    `json:"requestId"`
	Code      int                       `json:"code"`
	Message   string                    `json:"message"`
	Detail    AudioStreamResponseDetail `json:"detail"`
}

type AudioStreamResponseDetail struct {
	Errorcode    int    `json:"errorcode`
	DupRequestId string `json:"dupRequestId"`
}

// 语音文件 检查同步返回数据
type VoiceFileResponse struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	RequestID string          `json:"requestId"`
	BtID      string          `json:"btId"`
	Detail    VoiceFileDetail `json:"detail"`
}
type VoiceFileDetail struct {
	AudioDetail   []map[string]any `json:"audioDetail"`
	AudioTags     map[string]any   `json:"audioTags"`
	AudioText     string           `json:"audioText"`
	AudioTime     int              `json:"audioTime"`
	Code          int              `json:"code"`
	RequestParams map[string]any   `json:"requestParams"`
	RiskLevel     string           `json:"riskLevel"`
}

// 流关闭回调
type CloseStreamResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

type TianWangResponse struct {
	Code      int            `json:"code"`
	Message   string         `json:"message"`
	RequestID string         `json:"requestId"`
	RiskLevel string         `json:"riskLevel"`
	Detail    map[string]any `json:"detail"`
}

type TianWangParams struct {
	EventId        string `json:"eventId"`
	TokenId        string `json:"tokenId"`
	Ip             string `json:"ip"`
	SmDeviceId     string `json:"smDeviceId"`
	Phone          string `json:"phone"`
	Channel        string `json:"channel"`
	Version        string `json:"version"`
	RegisterMethod string `json:"registerMethod"`
}
