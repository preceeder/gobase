package ginserver

import (
	"github.com/gin-gonic/gin"
	"github.com/preceeder/gobase/try"
	"github.com/preceeder/gobase/utils"
	"github.com/preceeder/gobase/utils/datetimeh"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type HttpError interface {
	GetCode() int // 正常情况都是 200， 错误情况一般是  403
	GetMap() map[string]any
	Error() string
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		requestId := strconv.FormatInt(datetimeh.Now().TimestampMilli(), 10) + utils.RandStr(3)
		c.Set("requestId", requestId)
		c.Next()

		cost := time.Since(start)
		slog.Info("",
			"method", c.Request.Method,
			"path", path,
			"requestId", requestId,
			"status", c.Writer.Status(),
			"query", query,
			"ip", c.ClientIP(),
			"errors", c.Errors.ByType(gin.ErrorTypePrivate),
			"cost", cost)

	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetString("requestId")
		defer try.CatchException(func(err any, trace string) {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var ResStatus int = 500
			if he, ok := err.(HttpError); ok {
				c.JSON(he.GetCode(), he.GetMap())
				ResStatus = 200
			}

			httpRequest, _ := httputil.DumpRequest(c.Request, true)

			if ResStatus == 500 {
				slog.Error("Recovery from panic ",
					"err", err,
					"request", httpRequest,
					"requestId", requestId,
				)
			}
			c.AbortWithStatus(ResStatus)

		})

		c.Next()
	}
}
