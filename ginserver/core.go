/*
File Name:  core.go
Description:  暂时没有用上
Author:      Chenghu
Date:       2023/10/8 10:06
Change Activity:
*/
package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/preceeder/gobase/utils"
	"log/slog"
	"net/http"
	"strconv"
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

// obj 必须是 指针
func (c GContext) GetBody(obj any) {
	err := c.Context.ShouldBindBodyWith(obj, binding.JSON)
	if err != nil {
		slog.Error("获取body 失败", "erorr", err.Error(), "requestId", c.RequestId)
		panic(BaseHttpError{Code: StatusCodeCommonErr, ErrorCode: CodeParameterError, Message: "Parameter abnormality"})
	}
}

func (c GContext) GetQuery(obj any) {
	err := c.Context.ShouldBindQuery(obj)
	if err != nil {
		slog.Error("获取query 失败", "erorr", err.Error(), "requestId", c.RequestId)
		panic(BaseHttpError{Code: StatusCodeCommonErr, ErrorCode: CodeParameterError, Message: "Parameter abnormality"})
	}
}

func (c GContext) QueryInt(key string) (int, error) {
	kd, err := strconv.Atoi(c.Query(key))
	return kd, err
}

func (c GContext) QueryInt64(key string) (int64, error) {
	kd, err := strconv.ParseInt(c.Query(key), 10, 64)
	return kd, err
}

func (c GContext) QueryPageSize() (int, int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", "20"))
	if err != nil {
		size = 20
	}
	return page, size
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
