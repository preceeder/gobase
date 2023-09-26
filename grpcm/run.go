/*
File Name:  run.py
Description:
Author:      Chenghu
Date:       2023/8/31 09:57
Change Activity:
*/
package grpcm

import (
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
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
			fmt.Println("other signal", s)
		}
	}
}

type RpcConfig struct {
	Addr string `json:"addr"`
	Name string `json:"name"`
}

func redaRpcConfig(c viper.Viper) (rc *RpcConfig) {
	rpcConfig := c.Sub("rpc")

	if rpcConfig == nil {
		fmt.Printf("rpc config is nil")
		os.Exit(1)
	}
	err := rpcConfig.Unmarshal(&rc)
	if err != nil {
		fmt.Printf("rpc config read error: " + err.Error())
		os.Exit(1)
	}
	return
}

//var rpcListner net.Listener

var rcpConfig *RpcConfig

func InitRpc(config viper.Viper) {
	rcpConfig = redaRpcConfig(config)
}

func Server(server *grpc.Server) {
	// 创建 Tcp 连接
	if rcpConfig == nil {
		slog.Error("rpc 还未初始化 请先调用 gobase.Init(rpc:true)")
	}
	rpcLi, err := net.Listen("tcp", rcpConfig.Addr)
	if err != nil {
		slog.Error("监听失败: %v", "error", err.Error())
	}
	slog.Info("开启监听： ", "addr", rcpConfig.Addr)
	//开启信号监听
	c := StartSignalLister()

	//开启信号处理
	go SignalHandler(c, func() {
		//平滑关闭
		server.GracefulStop()
	})

	//初始化 注册路由
	InitRpcRouter(server)
	// 在 gRPC 服务上注册反射服务
	reflection.Register(server)

	err = server.Serve(rpcLi)
	if err != nil {
		slog.Error("failed to serve: %v", err)
	}

}

func Client(interceptor ...grpc.UnaryClientInterceptor) *grpc.ClientConn {
	clt, err := grpc.Dial(rcpConfig.Addr, grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(interceptor...))
	if err != nil {
		slog.Error("连接 gPRC 服务失败", "error", err)
	}
	return clt
}
