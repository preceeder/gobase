/*
File Name:  try.py
Description:
Author:      Chenghu
Date:       2023/8/25 11:39
Change Activity:
*/
package try

import (
	"bytes"
	"fmt"
	"runtime"
	"slices"
	"strings"
)

func CatchException(handle func(err any, e string)) {
	if err := recover(); err != nil {
		e := printStackTrace(err)

		handle(err, e)
	}
}

var JumpPackage = []string{"try.CatchException", "gin.(*Context).Next", "gin.(*Engine).handleHTTPRequest",
	"gin.(*Engine).ServeHTTP", "runtime.goexit", "http.(*conn).serve", "http.serverHandler.ServeHTTP",
	"runtime.gopanic", "runtime.panicmem", "runtime.sigpanic"}

// 打印堆栈信息
func printStackTrace(err any) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v --> ", err)
	isu := false
	for i := 1; true; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
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
					fmt.Fprintf(buf, "%s:%d --> ", file, line)
				}

			}
		}

	}
	return buf.String()
}

func ttttt() {
	defer CatchException(func(err any, e string) {
		if err != nil {
			fmt.Println(e)
		}
	})
	panic("ss")
}
