package ginserver

import (
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/preceeder/gobase/try"
	"github.com/preceeder/gobase/utils"
	"github.com/preceeder/gobase/utils/datetimeh"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var DefaultHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Methods":     "OPTIONS,GET,POST,PUT,DELETE",
	"Access-Control-Allow-Headers":     "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,x-auth-version,x-auth-channel,x-auth-channel-detail,x-auth-package,x-auth-timestamp,x-auth-announce,x-auth-token,x-auth-app,x-auth-signature",
	"Access-Control-Expose-Headers":    "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type",
	"Access-Control-Allow-Credentials": "true",
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		for key, value := range DefaultHeaders {
			c.Header(key, value) // 可将将 * 替换为指定的域名`
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger(serverLogHide bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		requestId := strconv.FormatInt(datetimeh.Now().TimestampMilli(), 10) + utils.RandStr(3)
		c.Set("requestId", requestId)
		c.Next()

		cost := time.Since(start)
		if !serverLogHide {
			slog.Info("",
				"method", c.Request.Method,
				"path", path,
				"requestId", requestId,
				"userId", c.GetString("userId"),
				"status", c.Writer.Status(),
				"query", query,
				"ip", c.ClientIP(),
				"errors", c.Errors.ByType(gin.ErrorTypePrivate),
				"cost", cost)
		}
	}
}

// 发生错误后 外部的处理
var ginRecoveryMidFuncs = []func(c *gin.Context, code int, err any, trance string){}

func PushGinRecoveryMidFunc(fc ...func(c *gin.Context, code int, err any, trance string)) {
	ginRecoveryMidFuncs = append(ginRecoveryMidFuncs, fc...)
}

// GinRecovery recover掉项目可能出现的panic，并使用slog记录相关日志
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetString("requestId")
		defer try.CatchException(func(err any, trace string) {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var ResStatus int = 500
			if he, ok := err.(HttpError); ok {
				c.JSON(he.GetCode(), he.GetMap())
				ResStatus = he.GetCode()
			}
			params := GetRequestParams(c)
			slog.Error("Recovery from panic ",
				"err", err,
				"trace", trace,
				"request", params,
				"requestId", requestId,
			)

			c.AbortWithStatus(ResStatus)
			// 可以对不同的 code 做其他处理
			for _, f := range ginRecoveryMidFuncs {
				f(c, ResStatus, err, trace)
			}
		})
		c.Next()
	}
}

type ParamsData struct {
	Body  string
	Query url.Values
	Url   string
}

func (p ParamsData) String() string {
	str, _ := sonic.ConfigFastest.MarshalToString(p)
	return str
}

func GetRequestParams(c *gin.Context) ParamsData {

	var body []byte
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			body = cbb
		}
	}
	if body == nil {
		bo, err := io.ReadAll(c.Request.Body)
		if err != nil {
			body = []byte("")
		} else {
			body = bo
		}
	}

	query := c.Request.Form
	urlp := c.Request.RequestURI

	return ParamsData{Body: string(body), Query: query, Url: urlp}
}
