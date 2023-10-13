//   File Name:  signal.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/13 17:17
//    Change Activity:

package utils

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func StartSignalLister() chan os.Signal {
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT)
	return c
}

func SignalHandler(c chan os.Signal, f func()) {
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("Program Exit...", s)
			close(c)
			f()
		default:
			slog.Info("other signal", s)
		}
	}
}
