package ginserver

///*
//File Name:  core.py
//Description:  暂时没有用上
//Author:      Chenghu
//Date:       2023/10/8 10:06
//Change Activity:
//*/
//package ginserver
//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//type HandlerFunc func(c *Context)
//
//func Handle(h HandlerFunc) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx := &Context{
//			Context: c,
//		}
//		h(ctx)
//	}
//}
//
//type Context struct {
//	*gin.Context
//	RequestId string
//}
//
//func (c *Context) Success(data interface{}) {
//	c.JSON(http.StatusOK, gin.H{
//		"code": 200,
//		"msg":  "success",
//		"data": data,
//	})
//}
//
//func (c *Context) Fail(code int, errCode int, message string) {
//	c.JSON(code, gin.H{
//		"code": errCode,
//		"msg":  message,
//	})
//}
