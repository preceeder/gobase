/*
File Name:  headers.go
Description:
Author:      Chenghu
Date:       2023/10/12 10:23
Change Activity:
*/
package ginserver

type DefaultHeader struct {
	Announce      string `header:"x-auth-announce"`
	Channel       string `header:"x-auth-channel"`
	ChannelDetail string `header:"x-auth-channel-detail"`
	Timestamp     string `header:"x-auth-timestamp"`
	Version       string `header:"x-auth-version"`
	Package       string `header:"x-auth-package"`
	Token         string `header:"x-auth-token"`
	Signature     string `header:"x-auth-signature"`
}
