//   File Name:  requestClient.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/9 10:10
//    Change Activity:

package utils

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

func NewRequestClient() *resty.Client {
	var RequestClient = resty.New()
	RequestClient.SetHeaders(
		map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	)
	RequestClient.SetTimeout(3 * time.Second)
	RequestClient.SetTransport(&http.Transport{
		MaxIdleConnsPerHost:   100,              // 对于每个主机，保持最大空闲连接数为 10
		IdleConnTimeout:       30 * time.Second, // 空闲连接超时时间为 30 秒
		TLSHandshakeTimeout:   10 * time.Second, // TLS 握手超时时间为 10 秒
		ResponseHeaderTimeout: 20 * time.Second, // 等待响应头的超时时间为 20 秒
	})
	return RequestClient
}
