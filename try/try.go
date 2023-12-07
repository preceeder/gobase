/*
File Name:  try.go
Description:
Author:      Chenghu
Date:       2023/8/25 11:39
Change Activity:
*/
package try

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"
)

func CatchException(handle func(err any, trace string)) {
	if err := recover(); err != nil {
		trace := printStackTrace(err)
		handle(err, trace)
	}
}

var JumpPackage = []string{"try.CatchException", "gin.(*Context).Next", "gin.(*Engine).handleHTTPRequest",
	"gin.(*Engine).ServeHTTP", "runtime.goexit", "http.(*conn).serve", "http.serverHandler.ServeHTTP",
	"runtime.gopanic", "runtime.panicmem", "runtime.sigpanic"}

// 打印全部堆栈信息
func printStackTrace(err any) string {
	buf := new(bytes.Buffer)
	pwd, _ := os.Executable()
	fmt.Fprintf(buf, "%v --> ", err)
	isu := false
	for i := 1; true; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			file = strings.TrimPrefix(file, pwd)
			fmt.Fprintf(buf, "%s:%d", file, line)
			break
		} else {
			prevFunc := runtime.FuncForPC(pc).Name()
			if !isu {
				if strings.Contains(prevFunc, "try.CatchException") {
					isu = true
				}
			} else {
				names := strings.Split(prevFunc, "/")
				if !slices.Contains(JumpPackage, names[len(names)-1]) {
					file = strings.TrimPrefix(file, pwd)
					fmt.Fprintf(buf, "%s:%d --> ", file, line)
				}
			}
		}

	}
	return buf.String()
}

// 打印堆栈信息 指定调用栈的上一级信息
// funcName 函数名, step 从funcNmae开始记录多少层
func GetStackTrace(funcName string, step int) string {
	buf := new(bytes.Buffer)
	isu := false
	pwd, _ := os.Executable()
	for i := 1; true; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			file = strings.TrimPrefix(file, pwd)
			fmt.Fprintf(buf, "%s:%d", file, line)
			break
		} else {
			prevFunc := runtime.FuncForPC(pc).Name()
			if !isu {
				if strings.HasSuffix(prevFunc, funcName) {
					isu = true
				}
				continue
			}
			if step > 0 {
				file = strings.TrimPrefix(file, pwd)
				fmt.Fprintf(buf, "%s:%d ", file, line)
				if step > 1 {
					fmt.Fprintf(buf, " --> ")
				}
				step -= 1
				continue
			}
			break

		}
	}
	return buf.String()
}
