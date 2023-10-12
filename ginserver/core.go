/*
File Name:  core.py
Description:  暂时没有用上
Author:      Chenghu
Date:       2023/10/8 10:06
Change Activity:
*/
package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/preceeder/gobase/utils"
	"golang.org/x/net/context"
	"net/http"
)

type HandlerFunc func(c *GContext)

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &GContext{
			Context:   c,
			RequestId: c.GetString("requestId"),
			UContext: utils.Context{
				RequestId: c.GetString("requestId"),
			},
		}
		h(ctx)
	}
}

type GContext struct {
	*gin.Context
	RequestId string
	UContext  utils.Context
	UserId    string // 只有有token的时候才会存在
}

func (c *GContext) Success(data ...interface{}) {
	var res = gin.H{"code": 200, "success": true}
	if len(data) > 0 {
		res["data"] = data[0]
	}
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func (c *GContext) Fail(code int, errCode int, message string) {
	c.JSON(code, gin.H{
		"code": errCode,
		"msg":  message,
	})
	c.Abort()
}

func (c GContext) FailByOld(httpError HttpError) {
	c.JSON(httpError.GetCode(), httpError.GetMap())
	c.Abort()
}
